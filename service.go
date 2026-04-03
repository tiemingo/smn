package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/util"
)

type Header struct {
	Title   string   `yaml:"title"`
	Authors []string `yaml:"author"`
	Subject string   `yaml:"subtitle"`
}

// Header used for new notes
const header = `---
title: "%v"
subtitle: "%v"
author: [%v]
date: "\\today"
---

%v`

func createNote(path string) error {

	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := actualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}
	notePath := filepath.Join(notesDir, path+".md")

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
	defaultAuthors := []string{}
	for _, author := range cfg.DefaultAuthors {
		defaultAuthors = append(defaultAuthors, fmt.Sprintf("\"%v\"", author))
	}
	authors := strings.Join(defaultAuthors, ", ")
	subject := filepath.Base(filepath.Dir(path))
	if subject == "." {
		subject = ""
	}
	content := fmt.Sprintf(header, filepath.Base(path), subject, authors, template)
	if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create note(%v): %v", notePath, err)
	}

	// Open note and sync after
	return openNoteAndSync(cfg, notePath, true)
}

func editNote(path string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := actualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}
	notePath := filepath.Join(notesDir, path+".md")

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	return openNoteAndSync(cfg, notePath, false)
}

func removeNote(path string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := actualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}
	notePath := filepath.Join(notesDir, path+".md")

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	// Remove note
	if err := os.Remove(notePath); err != nil {
		return fmt.Errorf("failed to delete note: %v", err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg, fmt.Sprintf("delete note %v", path)); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}
	return nil
}

func buildNote(path string) error {
	if path == "" {
		return fmt.Errorf("you have to provide a note title")
	}

	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := actualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}
	notePath := filepath.Join(notesDir, path+".md")

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	fileName, err := buildFileName(cfg, notePath)
	if err != nil {
		return fmt.Errorf("failed to create file name: %v", err)
	}

	// Run build command
	replaceCommand := strings.NewReplacer("{note_path}", notePath, "{output_path}", filepath.Join(cfg.OutputDir, fileName))
	buildCommand := cfg.BuildCommand
	for i, commandElement := range buildCommand {
		replacedString, err := util.ReplaceWithHomeDir(replaceCommand.Replace(commandElement))
		if err != nil {
			return fmt.Errorf("failed to replace home directory: %v", err)
		}
		buildCommand = slices.Replace(buildCommand, i, i+1, replacedString)
	}

	if _, err, stderr := util.RunCommand(buildCommand[0], buildCommand[1:]...); err != nil {
		return fmt.Errorf("failed to run %v: %v, stderr: %v", cfg.BuildCommand, err, stderr)
	}

	return nil
}

func latestNotes(amount int) error {
	cfg := config.GetConfig()

	// Get notes directory
	notesDir, err := actualNotesDir(cfg)
	if err != nil {
		return fmt.Errorf("failed get notes dir(%v): %v", notesDir, err)
	}

	// Sync if wanted
	if err := syncIfWanted(cfg); err != nil {
		log.Printf("failed to sync, proceeding anyways, if you want to terminate the program upon sync error you can change this in the config: %v", err)
	}

	notes, err := util.GetSortedMarkdownFiles(notesDir)
	if err != nil {
		return fmt.Errorf("failed to get notes: %v", err)
	}

	if amount > 0 && amount < len(notes) {
		notes = notes[:amount]
	}
	fmt.Println(strings.Join(notes, "\n"))
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
