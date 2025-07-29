package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chainguard-dev/clog"
	"github.com/imjasonh/goapidiff/internal/differ"
	"github.com/imjasonh/goapidiff/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	repoPath string
	format   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goapidiff <old> <new>",
		Short: "Analyze Go API changes between two versions",
		Long: `goapidiff compares Go APIs between two versions and reports breaking changes.

The arguments can be either:
  - Git refs (branches, tags, commits): main feature-branch
  - Directories: old-version/ new-version/

Both arguments must be the same type (both git refs or both directories).`,
		Example: `  # Compare git branches
  goapidiff main feature-branch

  # Compare git tags
  goapidiff v1.0.0 v2.0.0

  # Compare directories
  goapidiff old-version/ new-version/`,
		Args: cobra.ExactArgs(2),
		RunE: run,
	}

	rootCmd.Flags().StringVar(&repoPath, "repo", ".", "repository path (default: current directory)")
	rootCmd.Flags().StringVar(&format, "format", "text", "output format: text, json, markdown")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	log := clog.FromContext(ctx)

	oldRef := args[0]
	newRef := args[1]

	// Auto-detect if arguments are directories or git refs
	oldIsDir := isDirectory(oldRef)
	newIsDir := isDirectory(newRef)

	if oldIsDir != newIsDir {
		return fmt.Errorf("both arguments must be directories or both must be git refs (old is directory: %v, new is directory: %v)", oldIsDir, newIsDir)
	}

	log.InfoContext(ctx, "analyzing API changes",
		"repo", repoPath,
		"old", oldRef,
		"new", newRef,
		"format", format,
		"dirs", oldIsDir,
	)

	d := differ.New(repoPath)
	var report *differ.Report
	var err error

	if oldIsDir {
		report, err = d.DiffDirs(ctx, oldRef, newRef)
	} else {
		report, err = d.Diff(ctx, oldRef, newRef)
	}
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}

	r := reporter.New(format)
	if err := r.Report(os.Stdout, report); err != nil {
		return fmt.Errorf("report generation failed: %w", err)
	}

	if report.HasBreakingChanges() {
		os.Exit(1)
	}
	return nil
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
