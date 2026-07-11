package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const sockPath = "/tmp/pulsed.sock"

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
)

func printLine(line string) {
	parts := strings.SplitN(line, " at ", 2)
	if len(parts) != 2 {
		fmt.Println(line)
		return
	}
	ip := parts[0]
	rest := parts[1]

	timeParts := strings.SplitN(rest, ": ", 2)
	if len(timeParts) != 2 {
		fmt.Println(line)
		return
	}
	ts := timeParts[0]
	text := timeParts[1]

	fmt.Printf("%s[%s]%s %s%s%s %s\n", colorGray, ts, colorReset, colorCyan, ip, colorReset, text)
}

func main() {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println("could not connect to service, is pulsed running?")
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println(colorGreen + "connected to pulse" + colorReset)
	fmt.Println(colorGray + "type a message and hit enter, ctrl+d to quit" + colorReset)
	fmt.Println()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			printLine(scanner.Text())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			fmt.Println(colorGray + "disconnected" + colorReset)
			return
		}
		conn.Write([]byte(line))
	}
}
