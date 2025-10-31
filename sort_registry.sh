#!/bin/bash

# sort_registry.sh
# Sorts add_plugins.json by stars (high to low), then alphabetically by name
# Outputs the sorted result to plugins.json

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INPUT_FILE="${SCRIPT_DIR}/add_plugins.json"
OUTPUT_FILE="${SCRIPT_DIR}/plugins.json"
TEMP_FILE="${SCRIPT_DIR}/plugins.json.tmp"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is not installed. Please install it first:"
    echo "  sudo apt-get install jq   # Debian/Ubuntu"
    echo "  brew install jq           # macOS"
    exit 1
fi

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: add_plugins.json not found in $SCRIPT_DIR"
    exit 1
fi

echo "Sorting add_plugins.json by stars (descending) and name (alphabetical)..."

# Sort the JSON array:
# 1. Primary sort: by stars descending (treating null/missing as 0)
# 2. Secondary sort: by name alphabetically (case-insensitive)
jq 'sort_by([-((.stars // 0)), (.name | ascii_downcase)])' "$INPUT_FILE" > "$TEMP_FILE"

# Replace output file with sorted version
mv "$TEMP_FILE" "$OUTPUT_FILE"

echo "✓ Sorting complete!"
echo "  Input:  $INPUT_FILE (unchanged)"
echo "  Output: $OUTPUT_FILE (sorted)"
echo ""
echo "Top 5 plugins by stars:"
jq -r '.[:5] | .[] | "\(.stars // 0) ⭐ - \(.name)"' "$OUTPUT_FILE"
