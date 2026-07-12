package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// configPath returns the path to the file that stores this client's name.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "pulse", "name.txt"), nil
}

// resolveClientName checks whether a saved name exists on disk.
// If it does, it is returned directly.
// If not, the daemon prompts the connecting client over the socket,
// reads the response, and persists it for future connections.
func resolveClientName(conn net.Conn, scanner *bufio.Scanner) (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}

	// Returning client — read name from disk.
	if data, err := os.ReadFile(path); err == nil {
		name := strings.TrimSpace(string(data))
		if name != "" {
			return name, nil
		}
	}

	// First-time client — ask for a name over the connection.
	if _, err := conn.Write([]byte("Enter your name: ")); err != nil {
		return "", fmt.Errorf("failed to prompt for name: %w", err)
	}

	if !scanner.Scan() {
		return "", fmt.Errorf("client disconnected before providing a name")
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		return "", fmt.Errorf("client provided an empty name")
	}

	// Persist the name so future connections skip the prompt.
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(name), 0644); err != nil {
		return "", fmt.Errorf("failed to save name: %w", err)
	}

	return name, nil
}

// ListenIPC starts the Unix socket IPC listener and accepts CLI connections.
// For each new client it resolves the client name, replays history,
// registers the client, and starts forwarding its input over UDP.
func ListenIPC(sockPath string, history *History, clients *Clients, sendConn *net.UDPConn) error {
	unixListener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("error creating unix socket: %w", err)
	}

	go func() {
		defer unixListener.Close()

		// Accept CLI connections forever.
		for {
			conn, err := unixListener.Accept()
			if err != nil {
				fmt.Println("error accepting CLI connection:", err)
				continue
			}

			go handleClient(conn, history, clients, sendConn)
		}
	}()

	return nil
}

// handleClient resolves the client name, replays history, registers the
// client, and forwards its input over UDP. Runs in its own goroutine.
func handleClient(conn net.Conn, history *History, clients *Clients, sendConn *net.UDPConn) {
	scanner := bufio.NewScanner(conn)

	name, err := resolveClientName(conn, scanner)
	if err != nil {
		fmt.Println("error resolving client name:", err)
		conn.Close()
		return
	}

	// Send existing history to the newly connected client.
	for _, line := range history.All() {
		if _, err := conn.Write([]byte(line + "\n")); err != nil {
			fmt.Println("error sending history:", err)
			conn.Close()
			return
		}
	}

	clients.Add(conn, name)

	sendDataOverUDP(scanner, sendConn)

	clients.Remove(conn)
}

// sendDataOverUDP reads lines from a connected CLI and forwards them over UDP.
func sendDataOverUDP(scanner *bufio.Scanner, sendConn *net.UDPConn) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if _, err := sendConn.Write([]byte(line)); err != nil {
			fmt.Println("error sending UDP packet:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from client:", err)
	}
}
