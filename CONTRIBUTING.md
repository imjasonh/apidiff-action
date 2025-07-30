# Contributing

## Building the Action

Before committing changes that affect the action code, you must rebuild the dist folder:

```bash
npm run build
```

This compiles all dependencies into `dist/index.js` which is what GitHub Actions actually runs.

**Important**: The dist folder must be committed to the repository for the action to work.

## Node Version

This project requires Node.js v20. We recommend using nvm to manage Node versions:

```bash
nvm use
npm ci
npm run build
```

## Testing

Run tests with:

```bash
npm test
```

## Pre-commit Hooks

This project uses pre-commit hooks to ensure code quality. Install them with:

```bash
pre-commit install
```
