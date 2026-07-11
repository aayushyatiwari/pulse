# pulse

A tiny LAN chat system for a group on the same network (like a shared hotspot). No server, no accounts — everyone's machine talks directly to everyone else's over UDP broadcast.

Two pieces:

- **pulsed** — a background service. Always running. Listens for messages, keeps recent history in memory, sends desktop notifications when you're not actively chatting.
- **pulse** — the chat window. Run it, type, see everyone's messages live.

## How it works

`pulsed` runs once per machine, in the background, all the time. When you want to chat, you open a terminal and run `pulse` — it connects to your own local `pulsed`, shows you recent history, and streams new messages live. Close it whenever, `pulsed` keeps running and you won't miss anything sent while you were away (as long as `pulsed` itself hasn't restarted).

Messages go out over UDP broadcast, so everyone on the same network/hotspot sees them — no central server, no internet required.

## Install

```bash
git clone git@github.com:aayushyatiwari/pulse.git
cd pulse
go build -o pulsed .
cd cli && go build -o pulse . && cd ..
./install.sh
```

Works on any systemd-based Linux distro (Ubuntu, Arch, etc.) — no other dependencies.

## Usage

Once installed, from any terminal:

```bash
pulse
```

Type a message, hit enter. Ctrl+D to quit — `pulsed` keeps running in the background regardless.

## Notes

- History is in-memory only and resets if `pulsed` restarts.
- Desktop notifications fire only when no `pulse` window is currently open, so you don't get double-pinged.
- Requires a `notify-send`-compatible desktop notification daemon (ships by default on GNOME, KDE, most desktop Linux setups).

## Uninstall

```bash
sudo systemctl disable --now pulsed@$USER
sudo rm /usr/local/bin/pulsed /usr/local/bin/pulse /etc/systemd/system/pulsed@.service
sudo systemctl daemon-reload
```
