# apidiff-action

A GitHub Action that detects breaking API changes in Go code and comments on pull requests. This action uses the official [`golang.org/x/exp/cmd/apidiff`](https://pkg.go.dev/golang.org/x/exp/cmd/apidiff) tool to analyze API compatibility.

## Features

- 🔍 Automatically detects breaking API changes in pull requests
- 💬 Comments on PRs with detailed change reports
- ✅ Identifies both breaking and compatible changes
- 🚫 Can fail CI builds when breaking changes are detected
- 📊 Provides summary statistics in action outputs

## Usage

Add this action to your workflow:

```yaml
name: API Compatibility Check
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  apidiff:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write # Required for commenting on PRs
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Required to access base commit

      - uses: actions/setup-go@v5
        with:
          go-version: 'stable'

      - uses: imjasonh/apidiff-action@v1
        with:
          fail-on-breaking: true
```

## Inputs

| Input               | Description                                                 | Default               |
| ------------------- | ----------------------------------------------------------- | --------------------- |
| `working-directory` | Directory to run apidiff in                                 | `.`                   |
| `fail-on-breaking`  | Whether to fail the action if breaking changes are detected | `true`                |
| `comment-on-pr`     | Whether to comment on the PR with results                   | `true`                |
| `token`             | GitHub token for commenting on PRs                          | `${{ github.token }}` |

## Outputs

| Output                 | Description                                             |
| ---------------------- | ------------------------------------------------------- |
| `has-breaking-changes` | Whether breaking changes were detected (`true`/`false`) |
| `breaking-count`       | Number of breaking changes found                        |
| `compatible-count`     | Number of compatible changes found                      |

## Example PR Comment

The action will comment on your PR with a formatted report:

> # API Compatibility Check Results
>
> ## Summary
>
> | Type               | Count |
> | ------------------ | ----- |
> | Breaking changes   | 1     |
> | Compatible changes | 3     |
>
> ⚠️ **This PR contains breaking API changes!**
>
> ## Details
>
> ### `github.com/example/mypackage`
>
> #### ❌ Breaking changes
>
> - (\*Client).DoSomething: changed from func(string) error to func(context.Context, string) error
>
> #### ✅ Compatible changes
>
> - NewOption: added
> - WithTimeout: added
> - DefaultTimeout: added

## Advanced Usage

### Continue on breaking changes

To get notified about breaking changes without failing the build:

```yaml
- uses: imjasonh/apidiff-action@v1
  with:
    fail-on-breaking: false
```

### Check specific directory

To check a specific package or module:

```yaml
- uses: imjasonh/apidiff-action@v1
  with:
    working-directory: ./api
```

### Use outputs in subsequent steps

```yaml
- uses: imjasonh/apidiff-action@v1
  id: apidiff
  with:
    fail-on-breaking: false

- name: Handle breaking changes
  if: steps.apidiff.outputs.has-breaking-changes == 'true'
  run: |
    echo "Found ${{ steps.apidiff.outputs.breaking-count }} breaking changes"
    # Add custom handling here
```

## How it works

This action:

1. Installs the official `apidiff` tool from `golang.org/x/exp/cmd/apidiff`
2. Runs it to compare the base and head commits of your PR
3. Parses the output to identify breaking and compatible changes
4. Creates or updates a comment on the PR with the results
5. Optionally fails the build if breaking changes are detected

The underlying `apidiff` tool implements the compatibility rules from the [Go 1 compatibility guarantee](https://golang.org/doc/go1compat), checking for changes that would cause client code to stop compiling.

## License

This action is available under the [Apache License 2.0](LICENSE).
