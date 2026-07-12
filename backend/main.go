package main

import (
	"fmt"
	"os"
)

const sockPath = "/tmp/pulsed.sock" // Unix socket where IPC will happen

func main() {
	history := NewHistory(200)
	clients := NewClients()

	sendConn, err := DialUDP()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sendConn.Close()

	// Remove any stale Unix socket file from a previous run.
	_ = os.Remove(sockPath)

	if err := ListenIPC(sockPath, history, clients, sendConn); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := ListenUDP(history, clients); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	select {}
}
