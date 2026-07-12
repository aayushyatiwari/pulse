#!/usr/bin/env bash
# Create a GitHub release with cross-platform binaries.
# Usage: ./release.sh v1.0.0
set -e

VERSION=${1:?usage: ./release.sh <version>  e.g. ./release.sh v2.0.1}

TARGETS=(
    linux/amd64
    linux/arm64
    darwin/amd64
    darwin/arm64
    windows/amd64
)

# ── build ─────────────────────────────────────────────────────────────────────

rm -rf dist && mkdir dist

for TARGET in "${TARGETS[@]}"; do
    OS=${TARGET%/*}
    ARCH=${TARGET#*/}
    EXT=$([[ "$OS" == "windows" ]] && echo ".exe" || echo "")

    echo "→ building $OS/$ARCH"
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "dist/pulsed-$OS-$ARCH$EXT" ./backend/
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "dist/pulse-$OS-$ARCH$EXT"   ./frontend/
done

# ── release ───────────────────────────────────────────────────────────────────

echo "→ creating release $VERSION"
gh release create "$VERSION" dist/* \
    --title "Pulse $VERSION" \
    --notes "## Install

\`\`\`bash
bash <(curl -fsSL https://raw.githubusercontent.com/imaayush/pulse/main/install.sh)
\`\`\`

## Binaries

Download the binary pair for your platform from the assets below.
- \`pulsed\` — run on the hotspot host
- \`pulse\`  — run on each device to join the chat"

echo ""
echo "✓ release $VERSION published"
