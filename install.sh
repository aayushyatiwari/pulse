#!/bin/bash
set -e

systemctl --user stop pulsed 2>/dev/null || true

sudo cp pulsed /usr/local/bin/pulsed
sudo cp cli/pulse /usr/local/bin/pulse

mkdir -p ~/.config/systemd/user
cp pulsed.service ~/.config/systemd/user/pulsed.service

systemctl --user daemon-reload
systemctl --user enable --now pulsed

echo "installed. type 'pulse' in any terminal to chat."
echo ""
echo "note: pulsed normally stops when you log out."
echo "to keep it running even when logged out, run:"
echo "  loginctl enable-linger \$USER"
