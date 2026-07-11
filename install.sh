#!/bin/bash
set -e

sudo cp pulsed /usr/local/bin/pulsed
sudo cp cli/pulse /usr/local/bin/pulse
sudo cp pulsed.service /etc/systemd/system/pulsed@.service

sudo systemctl daemon-reload
sudo systemctl enable --now pulsed@$USER

echo "installed. type 'pulse' in any terminal to chat."
