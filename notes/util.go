package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiemingo/smn/note_config"
)

func (note *Note) getNotesDir() string {
	return note.notesDir
}

func (note *Note) getNoteRelDirParent() string {
	return filepath.Join(note.topic, note.subjects)
}

func (note *Note) getNoteDirParent() string {
	return filepath.Join(note.getNoteDir(), note.getNoteRelDirParent())
}

func (note *Note) getNoteRelDir() string {
	return filepath.Join(note.getNoteRelDirParent(), NameToNote(note.name))
}

func (note *Note) getNoteDir() string {
	return filepath.Join(note.getNotesDir(), note.getNoteRelDir())
}

func (note *Note) getEncryptedNotePath() string {
	return note.GetNotePath() + ".enc"
}

//
// Public non-static
//

func (note *Note) GetNotePath() string {
	return filepath.Join(note.getNoteDir(), "note.md")
}

func (note *Note) GetNoteRelName() string {
	return filepath.Join(note.getNoteRelDirParent(), note.name)
}

func (note *Note) GetTopicDir() string {
	return filepath.Join(note.getNotesDir(), note.topic)
}

func (note *Note) IsExist() bool {
	stat, err := os.Stat(note.GetNotePath())
	return err == nil && !stat.IsDir()
}

func (note *Note) GetConfig() note_config.Config {
	return note.config
}

//
// Private static
//

func splitRelNote(note string) (string, string, string, error) {

	parts := strings.Split(filepath.Clean(note), string(filepath.Separator))
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("provide atleast a topic and note name, subjects are optional")
	}

	topic := parts[0]
	parts = parts[1:]

	name := parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	subjects := strings.Join(parts, string(filepath.Separator))

	return name, topic, subjects, nil
}

func getExistingPathPart(base, rel string) string {
	current := filepath.Join(base, rel)

	for {
		// Check if the current path exists
		info, err := os.Stat(current)

		if err == nil && info.IsDir() {
			existingRel, err := filepath.Rel(base, current)
			if err != nil {
				return ""
			}
			return existingRel
		}

		// Get the parent directory
		parent := filepath.Dir(current)

		if parent == base {
			break
		}
		current = parent
	}

	return ""
}

//
// Public static
//

func NameToNote(name string) string {
	return fmt.Sprintf("note@%v", name)
}

func NoteToName(note string) string {
	return strings.TrimPrefix(note, "note@")
}

func IsNote(note string) bool {
	return strings.HasPrefix(strings.TrimSpace(note), "note@")
}
