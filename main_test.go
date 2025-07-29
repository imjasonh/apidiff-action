package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/imjasonh/goapidiff/internal/differ"
	"github.com/imjasonh/goapidiff/internal/reporter"
)

func TestDiffDirs(t *testing.T) {
	for _, tt := range []struct {
		name            string
		oldDir          string
		newDir          string
		wantBreaking    int
		wantCompatible  int
		wantExitCode    bool // true if should exit with code 1
		containsStrings []string
	}{{
		name:           "breaking changes",
		oldDir:         "testdata/break/old",
		newDir:         "testdata/break/new",
		wantBreaking:   1,
		wantCompatible: 4,
		wantExitCode:   true,
		containsStrings: []string{
			"(*Greeter).Greet: changed from func() string to func(bool) string",
			"Config.Timeout: added",
			"Greeter.Language: added",
			"Multiply: added",
			"StatusStarting: added",
		},
	}, {
		name:           "clean diff - no changes",
		oldDir:         "testdata/clean/old",
		newDir:         "testdata/clean/new",
		wantBreaking:   0,
		wantCompatible: 0,
		wantExitCode:   false,
		containsStrings: []string{
			"No API changes detected",
		},
	}, {
		name:           "safe diff - only additions",
		oldDir:         "testdata/safe/old",
		newDir:         "testdata/safe/new",
		wantBreaking:   0,
		wantCompatible: 13,
		wantExitCode:   false,
		containsStrings: []string{
			"(*Service).CallWithContext: added",
			"(*Service).GetStatus: added",
			"GetUsers: added",
			"UpdateUser: added",
			"User.CreatedAt: added",
			"User.Email: added",
			"StatusPending: added",
			"StatusUnknown: added",
		},
	}, {
		name:           "internal packages ignored",
		oldDir:         "testdata/internal/old",
		newDir:         "testdata/internal/new",
		wantBreaking:   1,
		wantCompatible: 1,
		wantExitCode:   true,
		containsStrings: []string{
			"PublicFunc: changed from func(string) string to func(string, bool) string",
			"PublicType.Email: added",
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			d := differ.New(".")

			report, err := d.DiffDirs(ctx, tt.oldDir, tt.newDir)
			if err != nil {
				t.Fatalf("DiffDirs failed: %v", err)
			}

			// Check counts
			gotBreaking := report.BreakingChangeCount()
			if gotBreaking != tt.wantBreaking {
				t.Errorf("BreakingChangeCount() = %d, want %d", gotBreaking, tt.wantBreaking)
			}

			gotCompatible := report.CompatibleChangeCount()
			if gotCompatible != tt.wantCompatible {
				t.Errorf("CompatibleChangeCount() = %d, want %d", gotCompatible, tt.wantCompatible)
			}

			// Check exit code behavior
			if report.HasBreakingChanges() != tt.wantExitCode {
				t.Errorf("HasBreakingChanges() = %v, want %v", report.HasBreakingChanges(), tt.wantExitCode)
			}

			// Check text output
			var buf bytes.Buffer
			r := reporter.New("text")
			if err := r.Report(&buf, report); err != nil {
				t.Fatalf("Report failed: %v", err)
			}

			output := buf.String()
			for _, want := range tt.containsStrings {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing expected string: %q\nGot output:\n%s", want, output)
				}
			}

			// For internal package test, ensure internal changes are NOT reported
			if tt.name == "internal packages ignored" {
				internalStrings := []string{"./internal", "Helper", "InternalType", "Process"}
				for _, shouldNotContain := range internalStrings {
					if strings.Contains(output, shouldNotContain) {
						t.Errorf("Output should not contain internal package info: %q found in output", shouldNotContain)
					}
				}
			}
		})
	}
}

func TestDiffDirsJSON(t *testing.T) {
	for _, tt := range []struct {
		name         string
		oldDir       string
		newDir       string
		wantBreaking bool
		packageCount int
	}{{
		name:         "breaking changes JSON",
		oldDir:       "testdata/break/old",
		newDir:       "testdata/break/new",
		wantBreaking: true,
		packageCount: 1,
	}, {
		name:         "clean diff JSON",
		oldDir:       "testdata/clean/old",
		newDir:       "testdata/clean/new",
		wantBreaking: false,
		packageCount: 0,
	}, {
		name:         "safe diff JSON",
		oldDir:       "testdata/safe/old",
		newDir:       "testdata/safe/new",
		wantBreaking: false,
		packageCount: 1,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			d := differ.New(".")

			report, err := d.DiffDirs(ctx, tt.oldDir, tt.newDir)
			if err != nil {
				t.Fatalf("DiffDirs failed: %v", err)
			}

			// Check JSON output
			var buf bytes.Buffer
			r := reporter.New("json")
			if err := r.Report(&buf, report); err != nil {
				t.Fatalf("Report failed: %v", err)
			}

			var result struct {
				OldRef             string        `json:"old_ref"`
				NewRef             string        `json:"new_ref"`
				HasBreakingChanges bool          `json:"has_breaking_changes"`
				BreakingCount      int           `json:"breaking_count"`
				CompatibleCount    int           `json:"compatible_count"`
				Packages           []interface{} `json:"packages"`
			}

			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if result.HasBreakingChanges != tt.wantBreaking {
				t.Errorf("JSON has_breaking_changes = %v, want %v", result.HasBreakingChanges, tt.wantBreaking)
			}

			if len(result.Packages) != tt.packageCount {
				t.Errorf("JSON packages count = %d, want %d", len(result.Packages), tt.packageCount)
			}
		})
	}
}

func TestDiffDirsMarkdown(t *testing.T) {
	ctx := context.Background()
	d := differ.New(".")

	// Test breaking changes with markdown format
	report, err := d.DiffDirs(ctx, "testdata/break/old", "testdata/break/new")
	if err != nil {
		t.Fatalf("DiffDirs failed: %v", err)
	}

	var buf bytes.Buffer
	r := reporter.New("markdown")
	if err := r.Report(&buf, report); err != nil {
		t.Fatalf("Report failed: %v", err)
	}

	output := buf.String()

	// Check markdown-specific formatting
	for _, want := range []string{
		"# API Changes Report",
		"**Old:** `testdata/break/old`",
		"**New:** `testdata/break/new`",
		"## Summary",
		"| Type | Count |",
		"| Breaking changes | 1 |",
		"| Compatible changes | 4 |",
		"### `./.`",
		"#### ❌ Breaking changes",
		"#### ✅ Compatible changes",
		"⚠️ **Warning:** This change contains breaking API changes!",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("Markdown output missing expected element: %q\nGot output:\n%s", want, output)
		}
	}
}

func TestIsDirectory(t *testing.T) {
	for _, tt := range []struct {
		path string
		want bool
	}{
		{"testdata/break/old", true},
		{"testdata/break/new", true},
		{"main.go", false},
		{"nonexistent", false},
		{"v1.2.3", false},
		{".", true},
	} {
		t.Run(tt.path, func(t *testing.T) {
			if got := isDirectory(tt.path); got != tt.want {
				t.Errorf("isDirectory(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
