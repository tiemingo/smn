package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/notes"
	"github.com/tiemingo/smn/util"
)

// syncIfWanted syncs the notes if the config option is set to auto sync.
// If an error occurs while syncing and the config option fail on sync error is enabled, the program exits with an error.
func syncIfWanted(cfg config.Config, optionalCommitMessage ...string) error {

	// Check if sync is wanted
	if !cfg.AutoSync {
		return nil
	}

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	// Loop through all topics to sync them
	repos, err := util.CollectGitRepos(notesDir)
	if err != nil {
		return fmt.Errorf("failed to collect git repos: %v", err)
	}
	for _, repo := range repos {
		if err := syncNotes(repo, optionalCommitMessage...); err != nil {

			// Check if should exit
			if cfg.FailOnSyncError {
				log.Fatalf("failed to sync, if you don't want the program to exit on failed sync, you can change it in the config: %v", err)
			}
			return err
		}
	}
	return nil
}

// openNoteAndSync opens the provided path with the default editor set in the env.
// The create bool decides whether the sync commit message should be create or update.
func openNoteAndSync(cfg config.Config, note *notes.Note, create bool) error {

	// Get editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, note.GetNotePath())
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
	syncIfWanted(cfg, fmt.Sprintf("%v note %v", mode, note.GetNoteRelName()))

	return nil
}

// GetSortedMarkdownFiles returns all notes names with topic and subjects in the provided and it's sub directories.
// The notes returned are ordered by last modified, with the most recent first.
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

		// Only process directories that have note@ prefix
		if d.IsDir() && notes.IsNote(filepath.Base(path)) {
			info, err := d.Info()
			if err != nil {
				return err
			}
			notePath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("failed to convert to relative path: %v", err)
			}
			notePath = filepath.Join(filepath.Dir(notePath), notes.NoteToName(filepath.Base(notePath)))
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
