#!/bin/bash
set -e

REPO="aayushyatiwari/pulse"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="$HOME/.config/systemd/user"

# --- detect architecture ---
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH_SLUG="amd64" ;;
  aarch64) ARCH_SLUG="arm64" ;;
  *)
    echo "error: unsupported architecture: $ARCH"
    echo "pulse currently supports x86_64 and aarch64 (arm64)."
    exit 1
    ;;
esac

# --- find latest release tag ---
echo "→ finding latest pulse release..."
TAG=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | head -1 \
  | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')

if [ -z "$TAG" ]; then
  echo "error: could not fetch latest release from GitHub."
  exit 1
fi

echo "→ installing pulse $TAG for linux/$ARCH_SLUG"

BASE_URL="https://github.com/$REPO/releases/download/$TAG"

# --- download binaries ---
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "→ downloading pulsed..."
curl -fsSL "$BASE_URL/pulsed-linux-$ARCH_SLUG" -o "$TMP/pulsed"

echo "→ downloading pulse..."
curl -fsSL "$BASE_URL/pulse-linux-$ARCH_SLUG" -o "$TMP/pulse"

chmod +x "$TMP/pulsed" "$TMP/pulse"

# --- download service file ---
echo "→ downloading service file..."
curl -fsSL "$BASE_URL/pulsed.service" -o "$TMP/pulsed.service"

# --- stop existing service if running ---
systemctl --user stop pulsed 2>/dev/null || true

# --- install binaries ---
echo "→ installing to $INSTALL_DIR (needs sudo)..."
sudo cp "$TMP/pulsed" "$INSTALL_DIR/pulsed"
sudo cp "$TMP/pulse"  "$INSTALL_DIR/pulse"

# --- install service ---
mkdir -p "$SERVICE_DIR"
cp "$TMP/pulsed.service" "$SERVICE_DIR/pulsed.service"

# --- enable and start ---
systemctl --user daemon-reload
systemctl --user enable --now pulsed

echo ""
echo "✓ pulse $TAG installed!"
echo ""
echo "  type 'pulse' in any terminal to start chatting."
echo ""
echo "note: pulsed stops when you log out. to keep it running even then:"
echo "  loginctl enable-linger \$USER"
