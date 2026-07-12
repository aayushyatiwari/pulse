#!/usr/bin/env bash
# curl -fsSL https://raw.githubusercontent.com/imaayush/pulse/main/install.sh | bash
set -e

REPO="imaayush/pulse"
BIN_DIR="$HOME/.local/bin"
SERVICE_DIR="$HOME/.config/systemd/user"

# ── platform detection ────────────────────────────────────────────────────────

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
    x86_64)          ARCH="amd64" ;;
    aarch64|arm64)   ARCH="arm64" ;;
    *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

if [[ "$OS" != "linux" && "$OS" != "darwin" ]]; then
    echo "unsupported OS: $OS"; exit 1
fi

# ── latest release ────────────────────────────────────────────────────────────

VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | cut -d'"' -f4)

echo "→ installing pulse $VERSION ($OS/$ARCH)"

# ── download binaries ─────────────────────────────────────────────────────────

mkdir -p "$BIN_DIR"

curl -fsSL "https://github.com/$REPO/releases/download/$VERSION/pulsed-$OS-$ARCH" \
    -o "$BIN_DIR/pulsed"
curl -fsSL "https://github.com/$REPO/releases/download/$VERSION/pulse-$OS-$ARCH" \
    -o "$BIN_DIR/pulse"

chmod +x "$BIN_DIR/pulsed" "$BIN_DIR/pulse"

# ensure ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$HOME/.bashrc"
    export PATH="$BIN_DIR:$PATH"
fi

# ── systemd user service (Linux) ──────────────────────────────────────────────

if [[ "$OS" == "linux" ]]; then
    mkdir -p "$SERVICE_DIR"

    cat > "$SERVICE_DIR/pulsed.service" << EOF
[Unit]
Description=Pulse LAN chat daemon
After=network.target

[Service]
ExecStart=$BIN_DIR/pulsed
Restart=on-failure
RestartSec=3

[Install]
WantedBy=default.target
EOF

    systemctl --user daemon-reload
    systemctl --user enable --now pulsed

    echo ""
    echo "✓ pulsed is running as a user service"

# ── launchd user agent (macOS) ────────────────────────────────────────────────

elif [[ "$OS" == "darwin" ]]; then
    PLIST_DIR="$HOME/Library/LaunchAgents"
    PLIST="$PLIST_DIR/com.pulse.pulsed.plist"
    mkdir -p "$PLIST_DIR"

    cat > "$PLIST" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>             <string>com.pulse.pulsed</string>
    <key>ProgramArguments</key>  <array><string>$BIN_DIR/pulsed</string></array>
    <key>RunAtLoad</key>         <true/>
    <key>KeepAlive</key>         <true/>
</dict>
</plist>
EOF

    launchctl unload "$PLIST" 2>/dev/null || true
    launchctl load -w "$PLIST"

    echo ""
    echo "✓ pulsed is running as a launch agent"
fi

# ── done ──────────────────────────────────────────────────────────────────────

echo "✓ pulse CLI installed"
echo ""
echo "  open the chat:     pulse"
echo "  daemon status:     systemctl --user status pulsed"
echo "  stop the daemon:   systemctl --user stop pulsed"
