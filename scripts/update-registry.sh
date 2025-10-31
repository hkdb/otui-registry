#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

# Create work directory
WORK_DIR="$REPO_ROOT/.work"
mkdir -p "$WORK_DIR"

echo "Fetching awesome-mcp-servers README..."
curl -o "$WORK_DIR/awesome-mcp-servers.md" \
  https://raw.githubusercontent.com/punkpeye/awesome-mcp-servers/main/README.md

echo "Parsing registry and enriching with GitHub metadata..."
if [ -z "$GITHUB_TOKEN" ]; then
    echo "WARNING: GITHUB_TOKEN not set. API rate limits may be hit."
    echo "Set GITHUB_TOKEN to increase rate limit from 60/hour to 5000/hour."
fi

cd "$REPO_ROOT"
go run scripts/parse-registry.go \
  "$WORK_DIR/awesome-mcp-servers.md" \
  "$WORK_DIR/registry.json"

# Move to repo root
mv "$WORK_DIR/results.json" "$REPO_ROOT/results.json"

echo "Done! Registry updated in results.json"
if command -v jq &> /dev/null; then
    echo "Plugins processed: $(jq length results.json)"
else
    echo "Install jq to see plugin count"
fi

# Cleanup work directory
rm -rf "$WORK_DIR"
