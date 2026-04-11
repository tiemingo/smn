package notes

import (
	"fmt"
	"os"
)

func (note *Note) Remove() error {
	if err := os.RemoveAll(note.getNoteDir()); err != nil {
		return fmt.Errorf("failed to delete note: %v", err)
	}
	return nil
}
