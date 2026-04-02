package main

import (
	"fmt"
	"os"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/util"
)

func addNote(path string) error {
	return nil
}

func editNote(path string) error {
	return nil
}

func buildNote(path string) error {
	return nil
}

func latestNotes(amount int) error {
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
	if statusOutput, statusErr, statusStderr := util.PipeInput("", "git", "status", "--porcelain"); statusErr != nil {
		return fmt.Errorf("failed to run git status --porcelain: %v, stderr: %v", statusErr, statusStderr)
	} else if statusOutput != "" {

		// Stage
		if _, addErr, addStderr := util.PipeInput("", "git", "add", "*"); addErr != nil {
			return fmt.Errorf("failed to run git add *: %v, stderr: %v", addErr, addStderr)
		}

		// Commit
		if _, commitErr, commitStderr := util.PipeInput("", "git", "commit", "-m", commitMessage); commitErr != nil {
			return fmt.Errorf("failed to run git commit -m \"%v\": %v, stderr: %v", commitMessage, commitErr, commitStderr)
		}
	}

	// Get updates from remote
	if _, pullErr, pullStderr := util.PipeInput("", "git", "pull", "--rebase"); pullErr != nil {
		return fmt.Errorf("failed to run git pull --rebase: %v, stderr: %v", pullErr, pullStderr)
	}

	// Push if changes exist
	if _, pushErr, pushStderr := util.PipeInput("", "git", "push"); pushErr != nil {
		return fmt.Errorf("failed to run git push: %v, stderr: %v", pushErr, pushStderr)
	}

	return nil
}
