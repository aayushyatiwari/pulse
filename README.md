# pulse

A tiny LAN chat for a group on the same network, like a shared hotspot. No server, no accounts, no internet needed — everyone's machine talks directly to everyone else's over UDP broadcast.

I made this to talk to my roommates when they connect to my hotspot. - Will not let them focus on their work. Lol.

- **pulsed** — a background service, always running, listens for messages, keeps recent history in memory, and sends a desktop notification when a message arrives and you're not actively chatting.
- **pulse** — the chat window itself. Run it, type, see everyone's messages live.

Type `pulse` in any terminal on the same network and you're in.

## Install

It's easy. Two things. One install the required packages, Two run the clone command. And you're done.

### The library code
#### Debian / Ubuntu 

```bash
sudo apt update
sudo apt install -y golang-go git libnotify-bin
```

#### Arch

```bash
sudo pacman -Sy --needed go git libnotify
```

### Build and install (same on every distro)

```bash
git clone https://github.com/aayushyatiwari/pulse.git
cd pulse
go build -o pulsed .
cd cli && go build -o pulse . && cd ..
./install.sh
```

Dont read the below expalantion. Just go and text people on your network! (people who have pulsed installed obviously)

This builds both binaries, copies them to `/usr/local/bin`, and sets up `pulsed` as a systemd **user service** — so it starts with your login session and can actually reach your desktop to show notifications.

By default, `pulsed` stops when you log out. To keep it running even while logged out (e.g. over SSH), run:

```bash
loginctl enable-linger $USER
```

## Usage

```bash
pulse
```

Type a message, hit enter. Ctrl+D to quit — `pulsed` keeps running in the background regardless, so you won't miss anything sent while `pulse` was closed (as long as `pulsed` itself hasn't restarted, which resets history).

## Checking it's running
To check if the background user service is running, run the below command.
```bash
systemctl --user status pulsed
```

Live logs:

```bash
journalctl --user -u pulsed -f
```

## Notes

- History is in-memory only and resets if `pulsed` restarts.
- Notifications fire only when no `pulse` window is open, so you don't get double-pinged.
- `pulsed` must run as a **user service**, not a system service — desktop notifications require access to your session's D-Bus bus, which system-level services don't have.

## Uninstall

```bash
systemctl --user disable --now pulsed
rm ~/.config/systemd/user/pulsed.service
sudo rm /usr/local/bin/pulsed /usr/local/bin/pulse
systemctl --user daemon-reload
```