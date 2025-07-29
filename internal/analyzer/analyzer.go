package analyzer

import (
	"context"
	"fmt"
	"go/types"
	"path/filepath"

	"github.com/chainguard-dev/clog"
	"golang.org/x/tools/go/packages"
)

type Analyzer struct {
	repoPath string
}

func New(repoPath string) *Analyzer {
	return &Analyzer{repoPath: repoPath}
}

func (a *Analyzer) LoadPackage(ctx context.Context, pkgPath string) (*types.Package, error) {
	log := clog.FromContext(ctx)

	absPath := filepath.Join(a.repoPath, pkgPath)
	log.InfoContext(ctx, "loading package", "path", absPath)

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir: a.repoPath,
	}

	// For directory comparison, we need to use the absolute path
	loadPath := pkgPath
	if filepath.IsAbs(a.repoPath) {
		loadPath = absPath
	}

	pkgs, err := packages.Load(cfg, loadPath)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found at %s", pkgPath)
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		log.WarnContext(ctx, "package has errors", "errors", pkg.Errors)
		// Continue anyway - apidiff can work with partially loaded packages
	}

	if pkg.Types == nil {
		return nil, fmt.Errorf("package %s has no type information", pkgPath)
	}

	return pkg.Types, nil
}

func (a *Analyzer) LoadPackages(ctx context.Context, pkgPaths []string) (map[string]*types.Package, error) {
	result := make(map[string]*types.Package)

	for _, pkgPath := range pkgPaths {
		pkg, err := a.LoadPackage(ctx, pkgPath)
		if err != nil {
			clog.FromContext(ctx).WarnContext(ctx, "failed to load package",
				"package", pkgPath,
				"error", err)
			continue
		}
		result[pkgPath] = pkg
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no packages could be loaded")
	}

	return result, nil
}
