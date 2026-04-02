package config

type Config struct {

	// Paths
	NotesDir   string `json:"notes_dir"`   // Directory that contains the notes. If sync is used this directory should contain the .git folder.
	BuildFiles string `json:"build_files"` // Directory that contains the build files for the notes.
	Template   string `json:"template"`    // File that contains a some default markdown.
	OutputDir  string `json:"output_dir"`  // Directory in which exported notes get saved.

	// Sync
	AutoSync        bool `json:"auto_sync"`          // If enabled, the program tries to sync the notes before and after all actions concerning them.
	FailOnSyncError bool `json:"fail_on_sync_error"` // If enabled the program will exit with an error if the sync fails.

	// Defaults
	DefaultAuthors string `json:"default_authors"` // Default authors used for the header in the markdown.

	// Build options
	BuildCommand     string `json:"build_command"`      // The command that gets run to export the markdown. The first %v gets replaced by the file name of the note and the second %v gets replaced by the output filename without filetype.
	BuildFileName    string `json:"build_file_name"`    // The name of the exported file. Replacers {authors} {subject} {title} can be used.
	BuildAuthor      string `json:"build_author"`       // How a single author should look in the file name.
	BuildAuthorSplit string `json:"build_author_split"` // How the authors should be chained together if multiple are present. Replacers {last_name} {middle_names} {fist_name} {given_name} can be used.
}
