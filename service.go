package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/notes"
	"github.com/tiemingo/smn/util"
)

func createNote(path string) error {

	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	note, err := notes.NoteObject(notesDir, path, true)
	if err != nil {
		return fmt.Errorf("failed to create note object: %v", err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	if err := note.CreateNote(); err != nil {
		return fmt.Errorf("failed to create note: %v", err)
	}

	// Open note and sync after
	return openNoteAndSync(cfg, note.GetNotePath(), true)
}

func editNote(path string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	note, err := notes.NoteObject(notesDir, path, false)
	if err != nil {
		return fmt.Errorf("failed to create note object: %v", err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	return openNoteAndSync(cfg, note.GetNotePath(), false)
}

func removeNote(path string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	// Remove note
	note, err := notes.NoteObject(notesDir, path, false)
	if err != nil {
		return fmt.Errorf("failed to create note object: %v", err)
	}
	if err := note.Remove(); err != nil {
		return err
	}

	// Sync if wanted
	if err := syncIfWanted(cfg, fmt.Sprintf("delete note %v", path)); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}
	return nil
}

func buildNote(path string, buildMode string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	note, err := notes.NoteObject(notesDir, path, false)
	if err != nil {
		return fmt.Errorf("failed to create note object: %v", err)
	}

	outputFilePath, err := note.BuildNote(buildMode)
	if err != nil {
		return fmt.Errorf("failed to build note: %v", err)
	}

	fmt.Println(outputFilePath)

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	return nil
}

func latestNotes(amount int) error {
	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := util.ReplaceWithHomeDir(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	notes, err := GetSortedMarkdownFiles(notesDir)
	if err != nil {
		return fmt.Errorf("failed to get notes: %v", err)
	}

	if amount > 0 && amount < len(notes) {
		notes = notes[:amount]
	}
	if len(notes) > 0 {
		fmt.Println(strings.Join(notes, "\n"))
	}
	return nil
}

func syncNotes(optionalCommitMessage ...string) error {

	// Check for optional commit message
	commitMessage := "update notes"
	if len(optionalCommitMessage) == 1 {
		commitMessage = optionalCommitMessage[0]
	}

	// Change wd to notes directory
	notesDir, err := util.ReplaceWithHomeDir(config.GetConfig().NotesDir)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	if err := os.Chdir(notesDir); err != nil {
		return fmt.Errorf("failed cd to notes dir(%v): %v", notesDir, err)
	}

	// Check for changes
	if statusOutput, statusErr, statusStderr := util.RunCommand("git", "status", "--porcelain"); statusErr != nil {
		return fmt.Errorf("failed to run git status --porcelain: %v, stderr: %v", statusErr, statusStderr)
	} else if statusOutput != "" {

		// Stage
		if _, addErr, addStderr := util.RunCommand("git", "add", "*"); addErr != nil {
			return fmt.Errorf("failed to run git add *: %v, stderr: %v", addErr, addStderr)
		}

		// Commit
		if _, commitErr, commitStderr := util.RunCommand("git", "commit", "-m", commitMessage); commitErr != nil {
			return fmt.Errorf("failed to run git commit -m \"%v\": %v, stderr: %v", commitMessage, commitErr, commitStderr)
		}
	}

	// Get updates from remote
	if _, pullErr, pullStderr := util.RunCommand("git", "pull", "--rebase"); pullErr != nil {
		return fmt.Errorf("failed to run git pull --rebase: %v, stderr: %v", pullErr, pullStderr)
	}

	// Check if something can be pushed to save time
	if cherryOutput, cherryErr, cherryStderr := util.RunCommand("git", "cherry", "-v"); cherryErr != nil {
		return fmt.Errorf("failed to run git cherry -v: %v, stderr: %v", cherryErr, cherryStderr)
	} else if cherryOutput != "" {

		// Push changes
		if _, pushErr, pushStderr := util.RunCommand("git", "push"); pushErr != nil {
			return fmt.Errorf("failed to run git push: %v, stderr: %v", pushErr, pushStderr)
		}
	}

	return nil
}
