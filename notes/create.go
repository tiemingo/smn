package notes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiemingo/smn/util"
)

// Header used for new notes
const noteHeader = `---
title: "%v"
subtitle: "%v"
author: [%v]
date: "\\today"
---

%v`

func (note *Note) CreateNote() error {

	noteDir := note.getNoteDir()

	// Check that note doesn't already exist
	if stat, err := os.Stat(noteDir); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed stat note(%v): %v", noteDir, err)
	} else if err == nil && stat.IsDir() {
		return fmt.Errorf("dir already exists at %v", noteDir)
	}

	// Create path to file
	if err := os.MkdirAll(noteDir, 0755); err != nil {
		return fmt.Errorf("failed to create subdirectories(%v): %v", noteDir, err)
	}

	// Load template
	template := ""
	if note.config.Template != "" {
		templatePath, err := util.ReplaceWithHomeDir(note.config.Template)
		if err != nil {
			return fmt.Errorf("failed convert template path(%v): %v", note.config.Template, err)
		}
		templateBytes, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to load template(%v): %v", templatePath, err)
		}
		template = string(templateBytes)
	}

	// Create note
	defaultAuthors := []string{}
	for _, author := range note.config.Authors {
		defaultAuthors = append(defaultAuthors, fmt.Sprintf("\"%v\"", author))
	}
	authors := strings.Join(defaultAuthors, ", ")
	subject := filepath.Base(note.subjects)
	if subject == "." {
		subject = ""
	}
	content := fmt.Sprintf(noteHeader, note.name, subject, authors, template)
	if err := os.WriteFile(note.GetNotePath(), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create note(%v): %v", noteDir, err)
	}

	return nil
}
