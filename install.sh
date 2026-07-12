#!/usr/bin/env bash
set -e

REPO="imaayush/pulse"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# ── checks ────────────────────────────────────────────────────────────────────

if ! command -v go &>/dev/null; then
    echo "error: Go is required — https://go.dev/dl/"
    exit 1
fi

# ── build ─────────────────────────────────────────────────────────────────────

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "→ cloning github.com/$REPO"
git clone --depth 1 "https://github.com/$REPO.git" "$TMP"

echo "→ building"
cd "$TMP"
go build -o pulsed ./backend/
go build -o pulse   ./frontend/

# ── install ───────────────────────────────────────────────────────────────────

echo "→ installing to $INSTALL_DIR (may require sudo)"
sudo install -m 755 pulsed pulse "$INSTALL_DIR"

echo ""
echo "✓ installed pulsed and pulse to $INSTALL_DIR"
echo ""
echo "  start the daemon:  pulsed"
echo "  open the chat:     pulse"
