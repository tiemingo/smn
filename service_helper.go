package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/util"
	"gopkg.in/yaml.v3"
)

// syncIfWanted syncs the notes if the config option is set to auto sync.
// If an error occurs while syncing and the config option fail on sync error is enabled, the program exits with an error.
func syncIfWanted(cfg config.Config, optionalCommitMessage ...string) error {

	// Check if sync is wanted
	if !cfg.AutoSync {
		return nil
	}

	err := syncNotes(optionalCommitMessage...)

	// Check if should exit
	if cfg.FailOnSyncError {
		log.Fatalf("failed to sync, if you don't want the program to exit on failed sync, you can change it in the config: %v", err)
	}

	return err
}

func openNoteAndSync(cfg config.Config, path string, create bool) error {

	// Get editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		syncIfWanted(cfg)
		return fmt.Errorf("failed to open file in editor(%v): %v", editor, err)
	}

	mode := "update"
	if create {
		mode = "create"
	}

	// Sync if wanted
	basePath, err := ActualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed to get notes directory: %v", err)
	}
	notePath, err := filepath.Rel(basePath, path)
	if err != nil {
		return fmt.Errorf("failed to convert to relative path: %v", err)
	}
	syncIfWanted(cfg, fmt.Sprintf("%v note %v", mode, notePath))

	return nil
}

func ActualNotesDir(cfg config.Config) (string, error) {
	basePath, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, "notes"), nil
}

func buildFileName(cfg config.Config, notePath string) (string, error) {

	header, err := getHeader(notePath)
	if err != nil {
		return "", fmt.Errorf("failed to get header: %v", err)
	}

	// Parse authors
	authors := []string{}
	for _, author := range header.Authors {
		replacedAuthor := util.ParseAuthor(author)
		authorReplacer := strings.NewReplacer("{last_name}", replacedAuthor.LastName, "{first_name}", replacedAuthor.FirstName, "{given_name}", replacedAuthor.GivenName)
		authors = append(authors, authorReplacer.Replace(author))
	}
	authorsString := strings.Join(authors, cfg.BuildAuthorSplit)

	titleReplacer := strings.NewReplacer("{authors}", authorsString, "{title}", header.Title, "{subject}", header.Subject)
	return titleReplacer.Replace(cfg.BuildFileName), nil
}

func getHeader(path string) (Header, error) {

	// Load note and extract header
	noteBytes, err := os.ReadFile(path)
	if err != nil {
		return Header{}, fmt.Errorf("failed to load note: %v", err)
	}
	parts := strings.SplitN(string(noteBytes), "---", 3)
	if len(parts) < 3 {
		return Header{}, fmt.Errorf("failed to find header: %v", err)
	}
	yamlBlock := parts[1]

	// Unmarshal
	var header Header
	err = yaml.Unmarshal([]byte(yamlBlock), &header)
	if err != nil {
		return Header{}, fmt.Errorf("failed to unmarshal header: %v", err)
	}

	return header, nil
}
