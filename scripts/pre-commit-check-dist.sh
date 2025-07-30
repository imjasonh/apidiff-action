#!/bin/bash
# Pre-commit script to check if dist is up to date

# Source nvm
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Use Node 20
nvm use 20.9.0 >/dev/null 2>&1

# Check if dist has uncommitted changes
if [ -n "$(git status --porcelain dist)" ]; then
  echo "dist is out of date. Run npm run build and commit the changes."
  exit 1
fi
