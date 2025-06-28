#!/bin/bash

# JAWT Installation Script
# Automatically detects system architecture, downloads the latest JAWT executable,
# installs it to /usr/local/bin, configures PATH, and sets permissions.

set -euo pipefail

# Configuration
REPO="yasufadhili/jawt"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="jawt"
TMP_DIR=$(mktemp -d)

# Ensure temporary directory is cleaned up on exit
trap 'rm -rf "$TMP_DIR"' EXIT

# Print error message and exit
error_exit() {
    echo "Error: $1" >&2
    exit 1
}

# Detect operating system and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$OS" in
    linux) OS_NAME="linux" ;;
    darwin) OS_NAME="macos" ;;
    *) error_exit "Unsupported operating system: $OS" ;;
esac

case "$ARCH" in
    x86_64) ARCH_NAME="amd64" ;;
    aarch64) ARCH_NAME="arm64" ;;
    *) error_exit "Unsupported architecture: $ARCH" ;;
esac

# Determine binary name based on OS and architecture
BINARY_PATTERN="jawt-${OS_NAME}-${ARCH_NAME}"
if [ "$OS" = "linux" ]; then
    BINARY_FILE="${BINARY_PATTERN}"
elif [ "$OS" = "darwin" ]; then
    BINARY_FILE="${BINARY_PATTERN}"
else
    error_exit "Unsupported platform: ${OS_NAME}/${ARCH_NAME}"
fi

# Check for required tools
command -v curl >/dev/null 2>&1 || error_exit "curl is required but not installed."
command -v jq >/dev/null 2>&1 || error_exit "jq is required but not installed."

# Fetch the latest release and find the appropriate asset
echo "Fetching latest JAWT release..."
DL_URL=$(curl -fsSL "$API_URL" | jq -r ".assets[] | select(.name | test(\"${BINARY_PATTERN}\")) | .browser_download_url")

if [ -z "$DL_URL" ]; then
    error_exit "No matching asset found for ${OS_NAME}/${ARCH_NAME}"
fi

# Download the binary
echo "Downloading JAWT binary from $DL_URL..."
curl -fsSL -o "$TMP_DIR/$BINARY_FILE" "$DL_URL" || error_exit "Failed to download JAWT binary"

# Install the binary
echo "Installing JAWT to $INSTALL_DIR..."
sudo mv "$TMP_DIR/$BINARY_FILE" "$INSTALL_DIR/$BINARY_NAME" || error_exit "Failed to install JAWT"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME" || error_exit "Failed to set executable permissions"

# Configure PATH
echo "Configuring PATH..."
SHELL_CONFIG_FILES=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile" "$HOME/.zsh_profile")
CONFIG_UPDATED=false
for CONFIG_FILE in "${SHELL_CONFIG_FILES[@]}"; do
    if [ -f "$CONFIG_FILE" ]; then
        # Check if PATH already includes INSTALL_DIR
        if ! grep -q "$INSTALL_DIR" "$CONFIG_FILE"; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$CONFIG_FILE"
            echo "Updated $CONFIG_FILE with PATH"
        fi
        # Source the file to apply changes in the current session
        source "$CONFIG_FILE" 2>/dev/null || true
        CONFIG_UPDATED=true
        break
    fi
done

if [ "$CONFIG_UPDATED" = false ]; then
    echo "Warning: No supported shell configuration file found (.bashrc, .zshrc, .bash_profile, .zsh_profile)."
    echo "Please add $INSTALL_DIR to your PATH manually."
fi

echo "JAWT has been successfully installed to $INSTALL_DIR/$BINARY_NAME"
echo "You can now run 'jawt' from the command line."