#!/bin/bash
# Seshmark wrapper for any CLI tool
#
# Usage:
#   1. Copy this script and rename it to match your tool
#   2. Replace TOOL_NAME with your tool's name
#   3. Make sure ~/bin comes first in PATH
#   4. Run your tool normally

TOOL_NAME="${TOOL_NAME:-my-tool}"
REAL_BIN=$(which -a "$TOOL_NAME" | grep -v "$HOME/bin" | head -1)

if [ -z "$REAL_BIN" ]; then
  echo "Error: Real $TOOL_NAME not found in PATH"
  exit 1
fi

SESSION_ID="${TOOL_NAME}:$(uuidgen 2>/dev/null || python3 -c 'import uuid; print(uuid.uuid4())' 2>/dev/null || date +%s)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="$TOOL_NAME"
export SESHMARK_MODEL="${SESHMARK_MODEL:-}"

exec "$REAL_BIN" "$@"
