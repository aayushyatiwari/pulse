#!/bin/bash
# Run this on your machine to cross-compile and publish a new GitHub Release.
# Usage: ./build_release.sh v1.2.0
#
# Requirements:
#   - Go installed
#   - gh (GitHub CLI) installed and authenticated: https://cli.github.com

set -e

TAG=${1:-}
if [ -z "$TAG" ]; then
  echo "usage: ./build_release.sh <tag>   (e.g. ./build_release.sh v1.0.0)"
  exit 1
fi

DIST="dist"
rm -rf "$DIST"
mkdir -p "$DIST"

echo "→ building for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o "$DIST/pulsed-linux-amd64" .
GOOS=linux GOARCH=amd64 go build -o "$DIST/pulse-linux-amd64"  ./cli

echo "→ building for linux/arm64..."
GOOS=linux GOARCH=arm64 go build -o "$DIST/pulsed-linux-arm64" .
GOOS=linux GOARCH=arm64 go build -o "$DIST/pulse-linux-arm64"  ./cli

echo "→ copying service file..."
cp pulsed.service "$DIST/pulsed.service"

echo "→ creating GitHub release $TAG..."
gh release create "$TAG" \
  "$DIST/pulsed-linux-amd64" \
  "$DIST/pulse-linux-amd64" \
  "$DIST/pulsed-linux-arm64" \
  "$DIST/pulse-linux-arm64" \
  "$DIST/pulsed.service" \
  --title "pulse $TAG" \
  --notes "Pre-built binaries for Linux (amd64 and arm64).

## Install

\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/aayushyatiwari/pulse/main/get.sh | bash
\`\`\`"

echo ""
echo "✓ release $TAG published!"
echo "  https://github.com/aayushyatiwari/pulse/releases/tag/$TAG"
