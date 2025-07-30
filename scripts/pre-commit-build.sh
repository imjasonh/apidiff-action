#!/bin/bash
# Pre-commit script to build dist with Node 20

# Check if we're already on Node 20
node_version=$(node --version)
if [[ "$node_version" =~ ^v20\. ]]; then
  # Already on Node 20, just build
  npm run build
  exit 0
fi

# Try to use nvm if available
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
  . "$NVM_DIR/nvm.sh"
  nvm install 20.9.0 >/dev/null 2>&1
  nvm use 20.9.0 >/dev/null 2>&1

  # Verify we switched to Node 20
  node_version=$(node --version)
  if [[ "$node_version" =~ ^v20\. ]]; then
    npm run build
    exit 0
  fi
fi

# If we're not on Node 20 and can't switch, error out
echo "Error: Node.js v20 is required (found $node_version)"
echo "Please install Node 20 or use nvm to switch versions"
exit 1
