#!/usr/bin/env bash
set -euo pipefail

REPO="adialaleal/odins"
BINARY="odins"
INSTALL_DIR="/usr/local/bin"

echo ""
echo "  ____  ____  ___ _   _ ____"
echo " / __ \\|  _ \\|_ _| \\ | / ___|"
echo "| |  | | | | || ||  \\| \\___ \\"
echo "| |__| | |_| || || |\\  |___) |"
echo " \\____/|____/|___|_| \\_|____/ "
echo ""
echo "  The All-Father of Local DNS — Installer"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" != "darwin" ]; then
  echo "  ✗ ODINS currently only supports macOS."
  exit 1
fi

case "$ARCH" in
  arm64)  ARCH_LABEL="arm64" ;;
  x86_64) ARCH_LABEL="amd64" ;;
  *)
    echo "  ✗ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Fetch latest release tag
echo "  → Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "  ✗ Could not determine latest release. Check https://github.com/${REPO}/releases"
  exit 1
fi

echo "  → Latest version: ${LATEST}"

# Build download URL
TARBALL="odins_${LATEST}_darwin_${ARCH_LABEL}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

# Download
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "  → Downloading ${TARBALL}..."
curl -fsSL "$URL" -o "$TMP_DIR/${TARBALL}"

# Extract
tar -xzf "$TMP_DIR/${TARBALL}" -C "$TMP_DIR"

# Install
echo "  → Installing to ${INSTALL_DIR}/${BINARY}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  chmod +x "${INSTALL_DIR}/${BINARY}"
else
  sudo mv "$TMP_DIR/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY}"
fi

echo ""
echo "  ✓ ODINS ${LATEST} installed at ${INSTALL_DIR}/${BINARY}"
echo ""
echo "  Run 'odins init' to get started!"
echo ""
