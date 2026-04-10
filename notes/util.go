package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (note *Note) getNotesDir() string {
	return note.notesDir
}

func (note *Note) getNoteRelDirParent() string {
	return filepath.Join(note.topic, note.subjects)
}

func (note *Note) getNoteRelDir() string {
	return filepath.Join(note.getNoteRelDirParent(), nameToNote(note.name))
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

//
// Private static
//

func nameToNote(name string) string {
	return fmt.Sprintf("note@%v", name)
}

func noteToName(note string) string {
	return strings.TrimPrefix("note@", note)
}

func isNote(note string) bool {
	return strings.HasPrefix("note@", note)
}

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
