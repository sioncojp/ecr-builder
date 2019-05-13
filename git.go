package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitGetCommitHash ... get commit hash of specified branch
func GitGetCommitHash() error {
	out, err := exec.Command(
		"git",
		"rev-parse",
		"HEAD",
	).Output()

	if err != nil {
		return fmt.Errorf("GitGetCommitHash: %s", err)
	}

	CommitHash = strings.TrimRight(string(out), "\n")
	return nil
}
