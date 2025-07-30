#!/bin/bash
# Pre-commit script to build dist with Node 20

# Source nvm
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
  . "$NVM_DIR/nvm.sh"
else
  echo "Error: nvm not found at $NVM_DIR/nvm.sh"
  exit 1
fi

# Install and use Node 20
nvm install 20.9.0 >/dev/null 2>&1
nvm use 20.9.0 >/dev/null 2>&1

# Verify Node version
node_version=$(node --version)
if [[ ! "$node_version" =~ ^v20\. ]]; then
  echo "Error: Failed to switch to Node 20 (current: $node_version)"
  exit 1
fi

# Build dist
npm run build
