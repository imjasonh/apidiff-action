# goapidiff

A Go tool that analyzes API changes between two versions of Go code and reports breaking changes. It can compare both git refs (branches, tags, commits) and directories.

## Features

- 🔍 Detects breaking API changes in Go code
- ✅ Identifies safe, compatible additions
- 📁 Compare directories or git refs (branches/tags/commits)
- 🚫 Ignores changes in `internal/` packages by default
- 📊 Multiple output formats: text, JSON, markdown
- 🚀 CI-friendly: exits with code 1 when breaking changes are detected

## Installation

```bash
go install github.com/imjasonh/goapidiff@latest
```

## Usage

```bash
# Compare git branches
goapidiff main feature-branch

# Compare git tags
goapidiff v1.0.0 v2.0.0

# Compare directories
goapidiff old-version/ new-version/

# Output as JSON
goapidiff --format json v1.0.0 v2.0.0

# Output as Markdown (great for PR comments)
goapidiff --format markdown main HEAD
```

## Examples

### Breaking Change Detection

```bash
$ goapidiff v1.0.0 v2.0.0

API Changes Report
==================
Old: v1.0.0
New: v2.0.0

Summary:
  Breaking changes:   1
  Compatible changes: 4

Package: ./
-----------

Breaking changes:
  ✗ (*Greeter).Greet: changed from func() string to func(bool) string

Compatible changes:
  ✓ Config.Timeout: added
  ✓ Greeter.Language: added
  ✓ Multiply: added
  ✓ StatusStarting: added

⚠️  This change contains breaking API changes!
```

### Clean Diff (No API Changes)

```bash
$ goapidiff v1.0.0 v1.0.1

API Changes Report
==================
Old: v1.0.0
New: v1.0.1

No API changes detected.
```

## How It Works

This tool leverages the excellent [`golang.org/x/exp/apidiff`](https://pkg.go.dev/golang.org/x/exp/apidiff) package, which does all the heavy lifting of analyzing Go APIs and determining compatibility.

The `apidiff` package implements the compatibility rules from the [Go 1 compatibility guarantee](https://golang.org/doc/go1compat), checking for changes that would cause client code to stop compiling.

`goapidiff` adds:
- Git integration for comparing refs
- Directory comparison support
- Multiple output formats
- Automatic filtering of `internal/` packages
- CI-friendly exit codes

## Internal Package Handling

By default, `goapidiff` ignores all changes in `internal/` packages since they are not part of the public API. According to Go conventions, internal packages cannot be imported by external packages, so breaking changes there don't affect API compatibility.

## Output Formats

### Text (default)
Human-readable format with clear sections for breaking and compatible changes.

### JSON
Machine-readable format, perfect for CI integration:

```json
{
  "old_ref": "v1.0.0",
  "new_ref": "v2.0.0",
  "has_breaking_changes": true,
  "breaking_count": 1,
  "compatible_count": 4,
  "packages": [
    {
      "package": "./",
      "changes": [
        {
          "message": "(*Greeter).Greet: changed from func() string to func(bool) string",
          "compatible": false
        }
      ]
    }
  ]
}
```

### Markdown
Great for posting as PR comments:

```markdown
# API Changes Report

**Old:** `v1.0.0`  
**New:** `v2.0.0`  

## Summary

| Type | Count |
|------|-------|
| Breaking changes | 1 |
| Compatible changes | 4 |

### `./`

#### ❌ Breaking changes

- (*Greeter).Greet: changed from func() string to func(bool) string

#### ✅ Compatible changes

- Config.Timeout: added
- Greeter.Language: added
```

## Exit Codes

- `0`: No breaking changes detected
- `1`: Breaking changes detected (useful for CI)

## Use in CI

Add to your GitHub Actions workflow:

```yaml
- name: Check API compatibility
  run: |
    go install github.com/imjasonh/goapidiff@latest
    goapidiff ${{ github.event.pull_request.base.sha }} ${{ github.event.pull_request.head.sha }}
```

## Credits

This tool is built on top of [`golang.org/x/exp/apidiff`](https://pkg.go.dev/golang.org/x/exp/apidiff), which provides the core API analysis functionality. All the smart API compatibility checking is done by that excellent package.


