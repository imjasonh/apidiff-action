#!/bin/bash
# Helper script to ensure Node 20 is being used

# Check if nvm is available
if ! command -v nvm &> /dev/null; then
  if [ -f "$HOME/.nvm/nvm.sh" ]; then
    source "$HOME/.nvm/nvm.sh"
  else
    echo "Error: nvm is not installed or not loaded"
    echo "Please install nvm or source it in your shell"
    exit 1
  fi
fi

# Install Node 20 if needed and switch to it
echo "Switching to Node 20..."
nvm install
nvm use

# Verify we're on Node 20
node_version=$(node --version)
if [[ ! "$node_version" =~ ^v20\. ]]; then
  echo "Error: Failed to switch to Node 20 (current: $node_version)"
  exit 1
fi

echo "✓ Using Node $node_version"
