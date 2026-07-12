package main

import (
	"fmt"
	"os/exec"
)

func Notify(line string) {
	err := exec.Command("notify-send", "pulse", line).Run()
	if err != nil {
		fmt.Printf("notify-send error: %v\n", err)
	}
}
