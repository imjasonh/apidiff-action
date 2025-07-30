# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a GitHub Action written in Node.js that detects breaking API changes in Go code using the official `golang.org/x/exp/cmd/apidiff` tool. It analyzes Go API compatibility between commits and can comment on pull requests with detailed change reports.

## Common Commands

### Development

**Important**: Use Node.js v20 to match GitHub Actions runtime.

```bash
# Install dependencies
npm ci

# Run tests
npm test

# Run linting
npm run lint

# Run formatter
npm run format

# Build the action (compiles to dist/)
npm run build

# Run pre-commit hooks manually
pre-commit run --all-files
```

### Testing the Action Locally

The action supports testing with environment variables:

```bash
# Test with directories
INPUT_OLD=testdata/break/old INPUT_NEW=testdata/break/new node index.js

# Test with git refs
INPUT_OLD=HEAD~1 INPUT_NEW=HEAD node index.js
```

## Architecture

### Core Components

1. **index.js** - Main entry point that:
   - Reads GitHub Action inputs
   - Determines which commits/directories to compare
   - Orchestrates the apidiff run, parsing, and PR commenting
   - Sets action outputs

2. **lib/apidiff.js** - Handles running the official apidiff tool:
   - Installs `golang.org/x/exp/cmd/apidiff` if not present
   - Supports both git refs and directory comparisons
   - For directories: creates export files first, then compares them
   - For git refs: runs apidiff directly on the commits

3. **lib/parser.js** - Parses apidiff text output:
   - Extracts breaking and compatible changes by package
   - Converts to structured format
   - Generates markdown reports for PR comments

4. **lib/commenter.js** - Manages PR comments:
   - Creates or updates existing comments (identified by hidden tag)
   - Requires `pull-requests: write` permission
   - Uses GitHub Actions toolkit for API access

### Testing Structure

The `testdata/` directory contains test cases:

- `clean/` - No API changes
- `safe/` - Only compatible changes (additions)
- `break/` - Contains breaking changes

Each test case has `old/` and `new/` subdirectories with Go code.

### Build Process

The action uses `@vercel/ncc` to compile all dependencies into a single `dist/index.js` file. This is required for GitHub Actions and must be committed to the repository. The dist/ folder is intentionally NOT in .gitignore.

**Important**: The dist must be built with Node.js v20 to match GitHub Actions runtime.

To build and commit dist after making changes:

```bash
# Method 1: Use the helper script
./scripts/use-node-20.sh  # Switches to Node 20
npm run build
git add dist/index.js
git commit -m "Rebuild dist"

# Method 2: Manual with nvm
nvm use  # Reads .nvmrc and switches to Node 20.9.0
node --version  # Should show v20.x.x
npm run build
git add dist/index.js
git commit -m "Rebuild dist"
```

**Important**: The pre-commit hooks will:
1. Verify you're using Node 20 before building dist
2. Automatically rebuild dist when source files change
3. Prevent commits if dist is out of date

If you get an error about Node version, switch to Node 20 using `nvm use` or `./scripts/use-node-20.sh`.

### Pre-commit Hooks

Configured with husky to run:

- ESLint and Prettier on JS files
- Pretty-format-json on JSON files (excluding package-lock.json)
- Builds dist/ and verifies it's up to date

## Key Design Decisions

1. **Node.js wrapper around Go tool**: Rather than reimplementing apidiff, this action wraps the official tool
2. **Support for both directories and git refs**: Allows testing without git commits
3. **Structured parsing**: Converts free-text apidiff output into structured data for better formatting
4. **Comment management**: Updates existing comments rather than creating duplicates

## Workflow Guidance

- When you open a PR, enable auto-merge, and watch CI checks until they complete and the PR merges. If CI fails, use the gh CLI to diagnose the error
