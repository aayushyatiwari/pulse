package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const sockPath = "/tmp/pulsed.sock"

func handleClient(conn net.Conn, sendConn *net.UDPConn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		sendConn.Write([]byte(line))
	}
}

func main() {
	history := NewHistory(200)
	clients := NewClients()

	broadcastAddr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	sendConn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		fmt.Println("dial error:", err)
		os.Exit(1)
	}
	defer sendConn.Close()

	os.Remove(sockPath)
	unixListener, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Println("unix listen error:", err)
		os.Exit(1)
	}
	defer unixListener.Close()

	go func() {
		for {
			conn, err := unixListener.Accept()
			if err != nil {
				fmt.Println("accept error:", err)
				continue
			}
			for _, line := range history.All() {
				conn.Write([]byte(line + "\n"))
			}
			clients.Add(conn)
			go handleClient(conn, sendConn)
		}
	}()

	listenAddr, _ := net.ResolveUDPAddr("udp4", ":9999")
	udpConn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, addr, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("read error:", err)
				continue
			}
			ip := addr.IP.String()
			ts := time.Now().Format("15:04:05")
			text := strings.TrimSpace(string(buf[:n]))
			line := fmt.Sprintf("%s at %s: \"%s\"", ip, ts, text)
			history.Add(line)
			clients.Broadcast(line)
			if clients.Count() == 0 {
				notify(line)
			}
			fmt.Println("[received]", line)
		}
	}()

	select {}
}
