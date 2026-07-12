package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const broadcastAddress = "255.255.255.255:9999"
const listenAddress = ":9999"

// DialUDP opens a UDP connection to the broadcast address and returns it.
func DialUDP() (*net.UDPConn, error) {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", broadcastAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve broadcast address: %w", err)
	}

	sendConn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		return nil, fmt.Errorf("dial error on UDP: %w", err)
	}

	return sendConn, nil
}

// ListenUDP starts the UDP listener and dispatches received pulse events
// to history and connected clients, triggering a desktop notification when
// no CLI clients are connected.
func ListenUDP(history *History, clients *Clients) error {
	listenAddr, err := net.ResolveUDPAddr("udp4", listenAddress)
	if err != nil {
		return fmt.Errorf("failed to resolve listen address: %w", err)
	}

	udpConn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return fmt.Errorf("listening error on UDP: %w", err)
	}

	go func() {
		defer udpConn.Close()
		buf := make([]byte, 1024)

		for {
			n, addr, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error reading UDP packet:", err)
				continue
			}

			raw := strings.TrimSpace(string(buf[:n]))

			if !strings.HasPrefix(raw, "PULSE:") {
				continue
			}

			payload := strings.TrimPrefix(raw, "PULSE:")
			parts := strings.SplitN(payload, "|", 2)
			if len(parts) != 2 {
				continue
			}

			name := parts[0]
			text := parts[1]
			ip := addr.IP.String()
			ts := time.Now().Format("15:04:05")

			sender := fmt.Sprintf("%s (%s)", name, ip)
			line := fmt.Sprintf("%s at %s: \"%s\"", sender, ts, text)

			history.Add(line)
			clients.Broadcast(line)

			if clients.Count() == 0 {
				Notify(line)
			}

			fmt.Println("[received]", line)
		}
	}()

	return nil
}
