package util

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Author struct {
	FirstName string
	LastName  string
	GivenName string
}

// RunCommand runs the provided command with args and returns the output, error and stderr
func RunCommand(command string, args ...string) (string, error, string) {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err, stderr.String()
	}

	return stdout.String(), nil, ""
}

// ReplaceWithHomeDir takes the path and replaces any occurrences of ~ with the home directory
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

// ParseAuthor parses the full name to first, last and give name
func ParseAuthor(fullName string) Author {

	parts := strings.Fields(fullName)

	if len(parts) == 0 {
		return Author{}
	}

	firstName := parts[0]
	lastName := ""
	givenNames := ""

	if len(parts) > 1 {
		lastName = parts[len(parts)-1]
		givenNames = strings.Join(parts[:len(parts)-1], " ")
	} else {
		givenNames = firstName
	}

	return Author{
		FirstName: firstName,
		LastName:  lastName,
		GivenName: givenNames,
	}
}

func CollectGitRepos(baseDir string) ([]string, error) {

	paths := []string{}

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			paths = append(paths, filepath.Dir(path))

			return filepath.SkipDir
		}
		return nil
	})
	return paths, err
}
