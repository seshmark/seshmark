package hook

const DelegatorScript = `#!/bin/sh
# Seshmark hook — delegates to seshmark binary
# Version: delegator-v1 (this script never changes)

# Allow bypass
[ -n "$SESHMARK_DISABLE" ] && exit 0

# Check if seshmark is available
if ! command -v seshmark >/dev/null 2>&1; then
  exit 0
fi

exec seshmark hook-run "$@"
`
