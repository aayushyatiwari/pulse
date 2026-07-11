package main

import (
	"os/exec"
	"fmt"
)
func notify(line string) {
	err := exec.Command("notify-send", "pulse", line).Run()
	if err != nil {
		fmt.Println("notify error:", err)
	}
}
