# Pulse

Pulse is a lightweight, LAN-based chat and event-broadcasting system. It allows users on the same network (e.g., a mobile hotspot) to chat seamlessly via UDP broadcast.

The project consists of:
- **`pulsed` (backend)**: A background daemon that listens for UDP broadcasts, maintains a history of recent messages, and dispatches desktop notifications when no chat client is open.
- **`pulse` (frontend)**: A beautiful terminal user interface (TUI) chat client that connects to your local `pulsed` daemon to read history, send messages, and view incoming chats in real-time.

## Prerequisites

For desktop notifications to work when the CLI is closed, the backend daemon (`pulsed`) relies on `notify-send`. You may need to install it depending on your Linux distribution:

- **Debian / Ubuntu / Linux Mint:**
  ```bash
  sudo apt install libnotify-bin
  ```
- **Arch Linux / Manjaro:**
  ```bash
  sudo pacman -S libnotify
  ```
- **Fedora:**
  ```bash
  sudo dnf install libnotify
  ```

*(Note: `notify-send` is currently Linux-specific. macOS users can still use the chat CLI perfectly fine, but background desktop notifications will silently skip if `notify-send` is unavailable.)*

## Installation

### Option 1: Automated Install (Recommended)

You can install both binaries and automatically set up the `pulsed` background service (via systemd on Linux or launchd on macOS) with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/aayushyatiwari/pulse/main/install.sh | bash
```

### Option 2: Download Pre-compiled Binaries

If you prefer not to run the install script, you can download the latest binaries directly from the [GitHub Releases](https://github.com/aayushyatiwari/pulse/releases) page.

1. Download the `pulse` and `pulsed` binaries for your platform (Linux/macOS, amd64/arm64).
2. Make them executable:
   ```bash
   chmod +x pulse pulsed
   ```
3. Move them to a directory in your PATH (e.g., `~/.local/bin` or `/usr/local/bin`):
   ```bash
   mv pulse pulsed ~/.local/bin/
   ```
4. Run `pulsed` in the background manually, or configure your own user service.

### Option 3: Build from Source (For Contributors)

If you want to contribute or build from source, you will need [Go 1.22+](https://go.dev/dl/) installed.

1. Clone the repository:
   ```bash
   git clone https://github.com/aayushyatiwari/pulse.git
   cd pulse
   ```
2. Build the binaries:
   ```bash
   go build -o pulsed ./backend/
   go build -o pulse ./frontend/
   ```
3. (Optional) Install them:
   ```bash
   sudo install -m 755 pulsed pulse /usr/local/bin/
   ```

## Usage

1. Make sure the backend daemon is running. (If you used the automated install script, it is already running as a background service).
   - *Manual start:* `./pulsed &`
2. Open the chat interface:
   ```bash
   pulse
   ```
3. The first time you run it, you'll be prompted to enter your display name.
4. Start chatting! Anyone on your local network (e.g., connected to the same hotspot) running `pulse` will receive your messages.

## Directory Structure

```text
pulse/
├── backend/          # All server-side daemon logic (pulsed)
│   ├── main.go       # Entry point for pulsed
│   ├── server.go     # Unix socket IPC server for local pulse clients
│   ├── udp.go        # UDP listener and broadcaster
│   ├── clients.go    # Manages connected local CLI clients
│   ├── history.go    # In-memory message history buffer
│   └── notify.go     # Desktop notification dispatcher
│
├── frontend/         # Terminal UI chat client (pulse)
│   └── main.go       # Entry point for the TUI client (using Bubble Tea)
│
├── install.sh        # Automated installation and service setup script
├── release.sh        # Script for building and publishing GitHub releases
└── README.md         # This file
```
