package differ

import (
	"context"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainguard-dev/clog"
	"github.com/imjasonh/goapidiff/internal/analyzer"
	"github.com/imjasonh/goapidiff/internal/git"
	"golang.org/x/exp/apidiff"
)

type Differ struct {
	repoPath string
}

func New(repoPath string) *Differ {
	return &Differ{repoPath: repoPath}
}

type Report struct {
	OldRef  string
	NewRef  string
	Changes []PackageChanges
}

type PackageChanges struct {
	Package string
	Report  apidiff.Report
}

func (r *Report) HasBreakingChanges() bool {
	for _, pc := range r.Changes {
		for _, change := range pc.Report.Changes {
			if !change.Compatible {
				return true
			}
		}
	}
	return false
}

func (r *Report) BreakingChangeCount() int {
	count := 0
	for _, pc := range r.Changes {
		for _, change := range pc.Report.Changes {
			if !change.Compatible {
				count++
			}
		}
	}
	return count
}

func (r *Report) CompatibleChangeCount() int {
	count := 0
	for _, pc := range r.Changes {
		for _, change := range pc.Report.Changes {
			if change.Compatible {
				count++
			}
		}
	}
	return count
}

func (d *Differ) Diff(ctx context.Context, oldRef, newRef string) (*Report, error) {
	log := clog.FromContext(ctx)

	repo, err := git.Open(d.repoPath)
	if err != nil {
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	changedPackages, err := repo.GetChangedGoPackages(ctx, oldRef, newRef)
	if err != nil {
		return nil, fmt.Errorf("getting changed packages: %w", err)
	}

	if len(changedPackages) == 0 {
		log.InfoContext(ctx, "no Go packages changed")
		return &Report{
			OldRef:  oldRef,
			NewRef:  newRef,
			Changes: []PackageChanges{},
		}, nil
	}

	// Create a temporary directory for git operations
	tmpDir, err := os.MkdirTemp("", "goapidiff-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Copy the repo to temp directory to avoid modifying the working directory
	if err := copyRepo(d.repoPath, tmpDir); err != nil {
		return nil, fmt.Errorf("copying repo: %w", err)
	}

	tmpRepo, err := git.Open(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("opening temp repository: %w", err)
	}

	// Load packages at old ref
	if err := tmpRepo.CheckoutRef(ctx, oldRef); err != nil {
		return nil, fmt.Errorf("checking out old ref: %w", err)
	}

	oldAnalyzer := analyzer.New(tmpDir)
	oldPackages, err := oldAnalyzer.LoadPackages(ctx, changedPackages)
	if err != nil {
		log.WarnContext(ctx, "failed to load old packages", "error", err)
		oldPackages = make(map[string]*types.Package)
	}

	// Load packages at new ref
	if err := tmpRepo.CheckoutRef(ctx, newRef); err != nil {
		return nil, fmt.Errorf("checking out new ref: %w", err)
	}

	newAnalyzer := analyzer.New(tmpDir)
	newPackages, err := newAnalyzer.LoadPackages(ctx, changedPackages)
	if err != nil {
		log.WarnContext(ctx, "failed to load new packages", "error", err)
		newPackages = make(map[string]*types.Package)
	}

	// Compare packages
	report := &Report{
		OldRef:  oldRef,
		NewRef:  newRef,
		Changes: []PackageChanges{},
	}

	for _, pkgPath := range changedPackages {
		oldPkg := oldPackages[pkgPath]
		newPkg := newPackages[pkgPath]

		if oldPkg == nil && newPkg == nil {
			continue
		}

		var changes apidiff.Report
		if oldPkg == nil {
			// New package
			changes = apidiff.Report{
				Changes: []apidiff.Change{{
					Message:    fmt.Sprintf("package %s added", pkgPath),
					Compatible: true,
				}},
			}
		} else if newPkg == nil {
			// Removed package
			changes = apidiff.Report{
				Changes: []apidiff.Change{{
					Message:    fmt.Sprintf("package %s removed", pkgPath),
					Compatible: false,
				}},
			}
		} else {
			// Compare packages
			changes = apidiff.Changes(oldPkg, newPkg)
		}

		if len(changes.Changes) > 0 {
			report.Changes = append(report.Changes, PackageChanges{
				Package: pkgPath,
				Report:  changes,
			})
		}
	}

	return report, nil
}

func (d *Differ) DiffDirs(ctx context.Context, oldDir, newDir string) (*Report, error) {
	log := clog.FromContext(ctx)

	// Make paths absolute
	absOldDir, err := filepath.Abs(oldDir)
	if err != nil {
		return nil, fmt.Errorf("getting absolute path for base dir: %w", err)
	}

	absNewDir, err := filepath.Abs(newDir)
	if err != nil {
		return nil, fmt.Errorf("getting absolute path for head dir: %w", err)
	}

	// Get all Go packages in both directories
	oldPackages, err := findGoPackages(absOldDir)
	if err != nil {
		return nil, fmt.Errorf("finding packages in base dir: %w", err)
	}

	newPackages, err := findGoPackages(absNewDir)
	if err != nil {
		return nil, fmt.Errorf("finding packages in head dir: %w", err)
	}

	// Combine all unique package paths
	allPackages := make(map[string]bool)
	for pkg := range oldPackages {
		allPackages[pkg] = true
	}
	for pkg := range newPackages {
		allPackages[pkg] = true
	}

	if len(allPackages) == 0 {
		log.InfoContext(ctx, "no Go packages found")
		return &Report{
			OldRef:  oldDir,
			NewRef:  newDir,
			Changes: []PackageChanges{},
		}, nil
	}

	// Load and compare packages
	oldAnalyzer := analyzer.New(absOldDir)
	newAnalyzer := analyzer.New(absNewDir)

	report := &Report{
		OldRef:  oldDir,
		NewRef:  newDir,
		Changes: []PackageChanges{},
	}

	for pkgPath := range allPackages {
		var oldPkg, newPkg *types.Package

		if oldPackages[pkgPath] {
			oldPkg, err = oldAnalyzer.LoadPackage(ctx, pkgPath)
			if err != nil {
				log.WarnContext(ctx, "failed to load base package", "package", pkgPath, "error", err)
			}
		}

		if newPackages[pkgPath] {
			newPkg, err = newAnalyzer.LoadPackage(ctx, pkgPath)
			if err != nil {
				log.WarnContext(ctx, "failed to load head package", "package", pkgPath, "error", err)
			}
		}

		if oldPkg == nil && newPkg == nil {
			continue
		}

		var changes apidiff.Report
		if oldPkg == nil {
			// New package
			changes = apidiff.Report{
				Changes: []apidiff.Change{{
					Message:    fmt.Sprintf("package %s added", pkgPath),
					Compatible: true,
				}},
			}
		} else if newPkg == nil {
			// Removed package
			changes = apidiff.Report{
				Changes: []apidiff.Change{{
					Message:    fmt.Sprintf("package %s removed", pkgPath),
					Compatible: false,
				}},
			}
		} else {
			// Compare packages
			changes = apidiff.Changes(oldPkg, newPkg)
		}

		if len(changes.Changes) > 0 {
			report.Changes = append(report.Changes, PackageChanges{
				Package: pkgPath,
				Report:  changes,
			})
		}
	}

	return report, nil
}

func findGoPackages(dir string) (map[string]bool, error) {
	packages := make(map[string]bool)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .go files (not test files)
		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			pkgPath := "./" + filepath.Dir(relPath)
			if pkgPath == "./" {
				pkgPath = "."
			}

			// Skip internal packages
			if strings.Contains(pkgPath, "/internal/") || strings.HasSuffix(pkgPath, "/internal") {
				return nil
			}

			packages[pkgPath] = true
		}

		return nil
	})

	return packages, err
}

func copyRepo(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
