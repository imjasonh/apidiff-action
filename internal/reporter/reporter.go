package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/imjasonh/goapidiff/internal/differ"
)

type Reporter struct {
	format string
}

func New(format string) *Reporter {
	return &Reporter{format: format}
}

func (r *Reporter) Report(w io.Writer, report *differ.Report) error {
	switch r.format {
	case "json":
		return r.reportJSON(w, report)
	case "markdown":
		return r.reportMarkdown(w, report)
	case "text":
		return r.reportText(w, report)
	default:
		return fmt.Errorf("unknown format: %s", r.format)
	}
}

func (r *Reporter) reportText(w io.Writer, report *differ.Report) error {
	fmt.Fprintf(w, "API Changes Report\n")
	fmt.Fprintf(w, "==================\n")
	fmt.Fprintf(w, "Old: %s\n", report.OldRef)
	fmt.Fprintf(w, "New: %s\n\n", report.NewRef)

	if len(report.Changes) == 0 {
		fmt.Fprintf(w, "No API changes detected.\n")
		return nil
	}

	breakingCount := report.BreakingChangeCount()
	compatibleCount := report.CompatibleChangeCount()

	fmt.Fprintf(w, "Summary:\n")
	fmt.Fprintf(w, "  Breaking changes:   %d\n", breakingCount)
	fmt.Fprintf(w, "  Compatible changes: %d\n\n", compatibleCount)

	for _, pc := range report.Changes {
		if len(pc.Report.Changes) == 0 {
			continue
		}

		fmt.Fprintf(w, "Package: %s\n", pc.Package)
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", len(pc.Package)+9))

		// Group changes by compatibility
		var breaking, compatible []string
		for _, change := range pc.Report.Changes {
			if change.Compatible {
				compatible = append(compatible, change.Message)
			} else {
				breaking = append(breaking, change.Message)
			}
		}

		if len(breaking) > 0 {
			fmt.Fprintf(w, "\nBreaking changes:\n")
			for _, msg := range breaking {
				fmt.Fprintf(w, "  ✗ %s\n", msg)
			}
		}

		if len(compatible) > 0 {
			fmt.Fprintf(w, "\nCompatible changes:\n")
			for _, msg := range compatible {
				fmt.Fprintf(w, "  ✓ %s\n", msg)
			}
		}

		fmt.Fprintf(w, "\n")
	}

	if breakingCount > 0 {
		fmt.Fprintf(w, "⚠️  This change contains breaking API changes!\n")
	}

	return nil
}

func (r *Reporter) reportJSON(w io.Writer, report *differ.Report) error {
	type jsonChange struct {
		Message    string `json:"message"`
		Compatible bool   `json:"compatible"`
	}

	type jsonPackageChanges struct {
		Package string       `json:"package"`
		Changes []jsonChange `json:"changes"`
	}

	type jsonReport struct {
		OldRef             string               `json:"old_ref"`
		NewRef             string               `json:"new_ref"`
		HasBreakingChanges bool                 `json:"has_breaking_changes"`
		BreakingCount      int                  `json:"breaking_count"`
		CompatibleCount    int                  `json:"compatible_count"`
		Packages           []jsonPackageChanges `json:"packages"`
	}

	jr := jsonReport{
		OldRef:             report.OldRef,
		NewRef:             report.NewRef,
		HasBreakingChanges: report.HasBreakingChanges(),
		BreakingCount:      report.BreakingChangeCount(),
		CompatibleCount:    report.CompatibleChangeCount(),
		Packages:           []jsonPackageChanges{},
	}

	for _, pc := range report.Changes {
		jpc := jsonPackageChanges{
			Package: pc.Package,
			Changes: []jsonChange{},
		}
		for _, change := range pc.Report.Changes {
			jpc.Changes = append(jpc.Changes, jsonChange{
				Message:    change.Message,
				Compatible: change.Compatible,
			})
		}
		jr.Packages = append(jr.Packages, jpc)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jr)
}

func (r *Reporter) reportMarkdown(w io.Writer, report *differ.Report) error {
	fmt.Fprintf(w, "# API Changes Report\n\n")
	fmt.Fprintf(w, "**Old:** `%s`  \n", report.OldRef)
	fmt.Fprintf(w, "**New:** `%s`  \n\n", report.NewRef)

	if len(report.Changes) == 0 {
		fmt.Fprintf(w, "*No API changes detected.*\n")
		return nil
	}

	breakingCount := report.BreakingChangeCount()
	compatibleCount := report.CompatibleChangeCount()

	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "| Type | Count |\n")
	fmt.Fprintf(w, "|------|-------|\n")
	fmt.Fprintf(w, "| Breaking changes | %d |\n", breakingCount)
	fmt.Fprintf(w, "| Compatible changes | %d |\n\n", compatibleCount)

	fmt.Fprintf(w, "## Changes by Package\n\n")

	for _, pc := range report.Changes {
		if len(pc.Report.Changes) == 0 {
			continue
		}

		fmt.Fprintf(w, "### `%s`\n\n", pc.Package)

		// Group changes by compatibility
		var breaking, compatible []string
		for _, change := range pc.Report.Changes {
			if change.Compatible {
				compatible = append(compatible, change.Message)
			} else {
				breaking = append(breaking, change.Message)
			}
		}

		if len(breaking) > 0 {
			fmt.Fprintf(w, "#### ❌ Breaking changes\n\n")
			for _, msg := range breaking {
				fmt.Fprintf(w, "- %s\n", msg)
			}
			fmt.Fprintf(w, "\n")
		}

		if len(compatible) > 0 {
			fmt.Fprintf(w, "#### ✅ Compatible changes\n\n")
			for _, msg := range compatible {
				fmt.Fprintf(w, "- %s\n", msg)
			}
			fmt.Fprintf(w, "\n")
		}
	}

	if breakingCount > 0 {
		fmt.Fprintf(w, "---\n\n")
		fmt.Fprintf(w, "⚠️ **Warning:** This change contains breaking API changes!\n")
	}

	return nil
}
