package main

import (
	"os"
	"os/exec"
	"strconv"
)

func FZFCommand() string {
	if TmuxPopupAvailable() {
		return "fzf-tmux -p"
	}
	return "fzf-tmux"
}

func TmuxPopupAvailable() bool {
	if _, found := os.LookupEnv("TMUX"); !found {
		return false
	}
	cmd := exec.Command("sh", "-c", "tmux -V | cut -d  ' ' -f2")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	versionStr := string(out)
	if len(versionStr) < 3 {
		return false
	}
	// remove version attribute such as 'a', 'b', 'rc-1' etc
	versionStr = versionStr[:3]
	version, err := strconv.ParseFloat(versionStr, 32)
	if err != nil {
		return false
	}
	if version < 3.2 {
		return false
	}
	return true
}
