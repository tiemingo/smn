package notes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiemingo/smn/note_config"
	"github.com/tiemingo/smn/util"
)

// Header used for new notes
const header = `---
title: "%v"
subtitle: "%v"
author: [%v]
date: "\\today"
---

%v`

type Note struct {
	notesDir string // Path to notes dir with replaced relatives
	topic    string // Name of the note topic. This is a direct child od the notes dir.
	subjects string // Relative path to the dir of the note, based on notesPath
	name     string // Name of the note

	aesGCM cipher.AEAD

	config note_config.Config
}

func NoteObject(notesDir, relNote string, new bool) (*Note, error) {

	name, topic, subjects, err := splitRelNote(relNote)
	if err != nil {
		return nil, fmt.Errorf("failed to split note path: %v", err)
	}

	note := &Note{
		notesDir: notesDir,
		subjects: subjects,
		name:     name,
		topic:    topic,
	}

	// Load config
	relDir := note.getNoteRelDir()
	if new {
		relDir = note.getNoteRelDirParent()

		// Go back until a folder exists
		relDir = getExistingPathPart(relDir)
		if relDir == "" {
			return nil, fmt.Errorf("make sure the topic folder exists")
		}
	}
	config, err := note_config.GetConfig(note.getNotesDir(), relDir)
	note.config = config

	// Set encryption if used
	var aesGCM cipher.AEAD
	if note.config.UseEncryption {
		if note.config.EncryptionKey == "" {
			return nil, fmt.Errorf("when using encryption please provide an encryption key")
		}

		block, err := aes.NewCipher([]byte(note.config.EncryptionKey))
		if err != nil {
			return nil, fmt.Errorf("failed to create cipher block: %v", err)
		}

		aes, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create aead: %v", err)
		}
		aesGCM = aes
	}
	note.aesGCM = aesGCM

	// If note is new create it
	if err := note.createNote(); err != nil {
		return nil, fmt.Errorf("failed to create note: %v", err)
	}

	return note, nil
}

func (note *Note) createNote() error {

	notePath := note.getNotePath()

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
	content := fmt.Sprintf(header, note.name, subject, authors, template)
	if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create note(%v): %v", notePath, err)
	}
	return nil
}
