package notes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/tiemingo/smn/note_config"
)

type Note struct {
	notesDir string // Path to notes dir with replaced relatives
	topic    string // Name of the note topic. This is a direct child od the notes dir.
	subjects string // Relative path to the dir of the note, based on notesPath
	name     string // Name of the note

	aesGCM cipher.AEAD

	config note_config.Config
}

func NoteObject(notesDir, relNote string, new bool, encryptionKey string) (*Note, error) {

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
		relDir = getExistingPathPart(note.getNotesDir(), relDir)
		if relDir == "" || relDir == "." {
			return nil, fmt.Errorf("make sure the topic folder exists")
		}
	}
	config, err := note_config.GetConfig(note.getNotesDir(), relDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get config for note: %v", err)
	}
	note.config = config

	// Set encryption if used
	var aesGCM cipher.AEAD
	if note.config.UseEncryption {
		if encryptionKey == "" {
			return nil, fmt.Errorf("when using encryption please provide an encryption key")
		}

		block, err := aes.NewCipher([]byte(encryptionKey))
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

	return note, nil
}
