package git

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chainguard-dev/clog"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct {
	repo *git.Repository
	path string
}

func Open(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("opening repository: %w", err)
	}
	return &Repository{repo: repo, path: path}, nil
}

func (r *Repository) GetChangedGoPackages(ctx context.Context, oldRef, newRef string) ([]string, error) {
	log := clog.FromContext(ctx)

	oldCommit, err := r.resolveCommit(oldRef)
	if err != nil {
		return nil, fmt.Errorf("resolving old ref %q: %w", oldRef, err)
	}

	newCommit, err := r.resolveCommit(newRef)
	if err != nil {
		return nil, fmt.Errorf("resolving new ref %q: %w", newRef, err)
	}

	log.InfoContext(ctx, "comparing commits",
		"old", oldCommit.Hash.String(),
		"new", newCommit.Hash.String(),
	)

	oldTree, err := oldCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("getting old tree: %w", err)
	}

	newTree, err := newCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("getting new tree: %w", err)
	}

	changes, err := oldTree.Diff(newTree)
	if err != nil {
		return nil, fmt.Errorf("diffing trees: %w", err)
	}

	packagesMap := make(map[string]bool)
	for _, change := range changes {
		from := change.From
		to := change.To

		if from.Name != "" && strings.HasSuffix(from.Name, ".go") && !strings.HasSuffix(from.Name, "_test.go") {
			pkg := getPackageFromPath(from.Name)
			if pkg != "" {
				packagesMap[pkg] = true
			}
		}
		if to.Name != "" && strings.HasSuffix(to.Name, ".go") && !strings.HasSuffix(to.Name, "_test.go") {
			pkg := getPackageFromPath(to.Name)
			if pkg != "" {
				packagesMap[pkg] = true
			}
		}
	}

	packages := make([]string, 0, len(packagesMap))
	for pkg := range packagesMap {
		// Skip internal packages
		if strings.Contains(pkg, "/internal/") || strings.HasSuffix(pkg, "/internal") {
			continue
		}
		packages = append(packages, pkg)
	}

	log.InfoContext(ctx, "found changed packages", "count", len(packages))
	return packages, nil
}

func (r *Repository) CheckoutRef(ctx context.Context, ref string) error {
	log := clog.FromContext(ctx)
	log.InfoContext(ctx, "checking out ref", "ref", ref)

	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	commit, err := r.resolveCommit(ref)
	if err != nil {
		return fmt.Errorf("resolving ref %q: %w", ref, err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash:  commit.Hash,
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("checking out commit: %w", err)
	}

	return nil
}

func (r *Repository) resolveCommit(ref string) (*object.Commit, error) {
	hash, err := r.repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return nil, fmt.Errorf("resolving revision: %w", err)
	}

	commit, err := r.repo.CommitObject(*hash)
	if err != nil {
		return nil, fmt.Errorf("getting commit object: %w", err)
	}

	return commit, nil
}

func getPackageFromPath(path string) string {
	dir := filepath.Dir(path)
	if dir == "." {
		return ""
	}
	return "./" + dir
}
