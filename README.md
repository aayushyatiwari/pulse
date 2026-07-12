# Pulse

Pulse is a lightweight notification and event-broadcasting system consisting of a background daemon (`pulsed`) and a command-line client (`pulse`).

---

## Directory Structure

```
pulse/
├── cmd/
│   ├── pulsed/          # Entry point for the pulsed daemon binary
│   └── pulse/           # Entry point for the pulse CLI client binary
│
├── internal/
│   ├── daemon/          # All server-side daemon logic
│   │   ├── server.go    # Core server: accepts IPC connections, coordinates subsystems
│   │   ├── history.go   # Event history: storage, retrieval, and pruning of past events
│   │   ├── clients.go   # Client registry: tracks connected pulse clients
│   │   ├── udp.go       # UDP listener: receives incoming pulse event packets
│   │   └── notify.go    # Notification dispatch: forwards events to registered clients
│   │
│   ├── client/          # Client-side library used by the pulse CLI
│   │   ├── ipc.go       # IPC transport: connects to and communicates with pulsed
│   │   ├── config.go    # Configuration: loads and parses client config files
│   │   └── register.go  # Registration: client handshake with the daemon
│   │
│   └── common/          # Shared types used by both daemon and client
│       └── protocol.go  # Wire protocol: message types, constants, and framing
│
└── README.md
```

## Binaries

| Binary  | Description                                                      |
|---------|------------------------------------------------------------------|
| `pulsed`| Background daemon. Listens for UDP events, manages clients, and dispatches notifications. |
| `pulse` | CLI client. Registers with `pulsed`, sends commands, and receives notifications.          |

## Building

```bash
# Build the daemon
go build ./cmd/pulsed

# Build the client
go build ./cmd/pulse
```

## Design Notes

- `internal/` packages are intentionally unexported and cannot be imported by external modules.
- `internal/common` is the only package shared between the daemon and client; it must remain free of platform-specific or subsystem-specific code.
- No third-party dependencies are used.
