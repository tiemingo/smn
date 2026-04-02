package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Returns the output, error and stderr
func PipeInput(input string, command string, args ...string) (string, error, string) {
	inputData := []byte(input)

	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err, stderr.String()
	}

	return stdout.String(), nil, ""
}

func ReplaceWithHomeDir(path string) (string, error) {
	if !strings.Contains(path, "~") {
		return path, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path, fmt.Errorf("failed to replace ~ with home directory: %v", err)
	}
	return strings.ReplaceAll(path, "~", homeDir), nil
}
