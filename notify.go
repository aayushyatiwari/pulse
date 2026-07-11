package main

import (
	"os/exec"
)

func notify(line string) {
	exec.Command("notify-send", "pulse", line).Run()
}
