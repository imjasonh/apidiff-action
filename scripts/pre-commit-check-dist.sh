#!/bin/bash
# Pre-commit script to check if dist is up to date

# Try to use Node 20 via nvm if available
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
  . "$NVM_DIR/nvm.sh"
  nvm use 20.9.0 >/dev/null 2>&1
fi

# Check if dist has uncommitted changes
if [ -n "$(git status --porcelain dist)" ]; then
  echo "dist is out of date. Run npm run build and commit the changes."
  exit 1
fi
