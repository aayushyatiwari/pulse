package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
	sender := parts[0]
	rest := parts[1]

	timeParts := strings.SplitN(rest, ": ", 2)
	if len(timeParts) != 2 {
		fmt.Println(line)
		return
	}
	ts := timeParts[0]
	text := timeParts[1]

	fmt.Printf("%s[%s]%s %s%s%s %s\n", colorGray, ts, colorReset, colorCyan, sender, colorReset, text)
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pulse")
}

func configPath() string {
	return filepath.Join(configDir(), "identity.conf")
}

// loadIdentity reads name and tag from the config file.
// Returns ("", "") if the file does not exist.
func loadIdentity() (name, tag string) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return "", ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name=") {
			name = strings.TrimPrefix(line, "name=")
		} else if strings.HasPrefix(line, "tag=") {
			tag = strings.TrimPrefix(line, "tag=")
		}
	}
	return
}

// saveIdentity writes name and tag to the config file.
func saveIdentity(name, tag string) error {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return err
	}
	content := fmt.Sprintf("name=%s\ntag=%s\n", name, tag)
	return os.WriteFile(configPath(), []byte(content), 0644)
}

// setupIdentity runs the first-time interactive setup and returns name and tag.
func setupIdentity(reader *bufio.Reader) (string, string) {
	fmt.Println(colorGreen + "welcome to pulse! let's set up your identity." + colorReset)
	fmt.Println()

	fmt.Print("your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "anonymous"
	}

	fmt.Print("your tag (e.g. 'family pandit', 'office dev', or leave blank): ")
	tag, _ := reader.ReadString('\n')
	tag = strings.TrimSpace(tag)

	if err := saveIdentity(name, tag); err != nil {
		fmt.Println(colorGray+"warning: could not save identity:", err, colorReset)
	} else {
		fmt.Println(colorGray + "identity saved to " + configPath() + colorReset)
	}
	fmt.Println()
	return name, tag
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	name, tag := loadIdentity()
	if name == "" {
		name, tag = setupIdentity(reader)
	}

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		fmt.Println("could not connect to service, is pulsed running?")
		os.Exit(1)
	}
	defer conn.Close()

	identityLine := name
	if tag != "" {
		identityLine = name + " [" + tag + "]"
	}
	fmt.Println(colorGreen+"connected to pulse as", identityLine+colorReset)
	fmt.Println(colorGray + "type a message and hit enter, ctrl+d to quit" + colorReset)
	fmt.Println()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			printLine(scanner.Text())
		}
	}()

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			fmt.Println(colorGray + "disconnected" + colorReset)
			return
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		// Wire format: PULSE:<name>|<tag>|<message>
		payload := fmt.Sprintf("PULSE:%s|%s|%s", name, tag, text)
		conn.Write([]byte(payload + "\n"))
	}
}
