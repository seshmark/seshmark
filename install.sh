#!/bin/bash
set -e

REPO="seshmark/seshmark"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;  # supported
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

BINARY="seshmark-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo "Installing seshmark..."
mkdir -p "$INSTALL_DIR"

curl -fsSL "$URL" -o "${INSTALL_DIR}/seshmark"
chmod +x "${INSTALL_DIR}/seshmark"

# Create aliases
ln -sf "${INSTALL_DIR}/seshmark" "${INSTALL_DIR}/agentblame"
ln -sf "${INSTALL_DIR}/seshmark" "${INSTALL_DIR}/aiblame"
ln -sf "${INSTALL_DIR}/seshmark" "${INSTALL_DIR}/git-agentblame"
ln -sf "${INSTALL_DIR}/seshmark" "${INSTALL_DIR}/git-aiblame"

# Ensure PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    for rc in ~/.bashrc ~/.zshrc; do
        if [ -f "$rc" ] && ! grep -q "$INSTALL_DIR" "$rc" 2>/dev/null; then
            echo "export PATH=\"${INSTALL_DIR}:\$PATH\"" >> "$rc"
        fi
    done
    echo "Added ${INSTALL_DIR} to PATH. Restart your shell or run:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
fi

# Install global hook
"${INSTALL_DIR}/seshmark" hook install --global

# Install in current repo if inside one
if git rev-parse --git-dir >/dev/null 2>&1; then
    "${INSTALL_DIR}/seshmark" hook install
fi

echo ""
echo "seshmark installed successfully!"
echo ""
echo "Try it now:"
echo "  git agentblame <file>"
echo "  seshmark query --agent cursor"
echo "  seshmark who HEAD"
echo ""
echo "To start tracking a session manually:"
echo "  seshmark track cursor:my-session"
echo ""
echo "To upgrade when a new version is released:"
echo "  seshmark upgrade"
