package config

type Config struct {

	// Paths

	NotesDir string `json:"notes_dir"` // Directory that contains the notes. If sync is used this directory should not contain the .git folder, it should be the parent of the directories with .git folder. But multiple git folde can be used

	// Sync

	AutoSync        bool `json:"auto_sync"`          // If enabled, the program tries to sync the notes before and after all actions concerning them.
	FailOnSyncError bool `json:"fail_on_sync_error"` // If enabled the program will exit with an error if the sync fails.
}
