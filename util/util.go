package util

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Author struct {
	FirstName string
	LastName  string
	GivenName string
}

// Returns the output, error and stderr
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

func GetSortedMarkdownFiles(root string) ([]string, error) {

	type FileInfo struct {
		Path    string
		ModTime int64
	}
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only process files ending in .md
		if !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			info, err := d.Info()
			if err != nil {
				return err
			}
			notePath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("failed to convert to relative path: %v", err)
			}
			files = append(files, FileInfo{
				Path:    notePath,
				ModTime: info.ModTime().Unix(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by by last modification
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime > files[j].ModTime
	})

	// Extract just the paths into a string slice
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.Path
	}

	return result, nil
}
