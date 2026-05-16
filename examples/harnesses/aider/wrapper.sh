#!/bin/bash
# Seshmark wrapper for Aider
#
# Install:
#   cp examples/harnesses/aider/wrapper.sh ~/bin/aider
#   chmod +x ~/bin/aider
#
# Make sure ~/bin comes before the real aider in PATH.

REAL_AIDER=$(which -a aider | grep -v "$HOME/bin" | head -1)

if [ -z "$REAL_AIDER" ]; then
  echo "Error: Real aider not found in PATH"
  exit 1
fi

# Generate session ID
SESSION_ID="aider:$(uuidgen 2>/dev/null || python3 -c 'import uuid; print(uuid.uuid4())' 2>/dev/null || date +%s)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="aider"

# Set model from aider's output if available, or let user configure
export SESHMARK_MODEL="${SESHMARK_MODEL:-}"

exec "$REAL_AIDER" "$@"
