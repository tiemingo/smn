package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/util"
)

const header = `---
title: "%v"
subject: "%v"
author: "%v"
date: "\\today"
---

%v`

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
	notePath := filepath.Join(notesDir, "notes", path+".md")

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	// Check that note doesn't already exist
	if stat, err := os.Stat(notePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed stat note(%v): %v", notePath, err)
	} else if err == nil && !stat.IsDir() {
		return fmt.Errorf("file already exists at %v", notePath)
	}

	// Create path to file
	if err := os.MkdirAll(filepath.Dir(notePath), 0755); err != nil {
		return fmt.Errorf("failed to create subdirectories(%v): %v", filepath.Dir(notePath), err)
	}

	// Load template
	template := ""
	if cfg.Template != "" {
		templatePath, err := util.ReplaceWithHomeDir(cfg.Template)
		if err != nil {
			return fmt.Errorf("failed convert template path(%v): %v", cfg.Template, err)
		}
		templateBytes, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to load template(%v): %v", cfg.Template, err)
		}
		template = string(templateBytes)
	}

	// Create note
	content := fmt.Sprintf(header, filepath.Base(path), filepath.Base(filepath.Dir(notePath)), cfg.DefaultAuthors, template)
	if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create note(%v): %v", notePath, err)
	}

	// Open note and sync after
	openNoteAndSync(cfg, notePath, true)

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

	// Check if something can be pushed to save time
	if cherryOutput, cherryErr, cherryStderr := util.PipeInput("", "git", "cherry", "-v"); cherryErr != nil {
		return fmt.Errorf("failed to run git cherry -v: %v, stderr: %v", cherryErr, cherryStderr)
	} else if cherryOutput != "" {

		// Push changes
		if _, pushErr, pushStderr := util.PipeInput("", "git", "push"); pushErr != nil {
			return fmt.Errorf("failed to run git push: %v, stderr: %v", pushErr, pushStderr)
		}
	}

	return nil
}
