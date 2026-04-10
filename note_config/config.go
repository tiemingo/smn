package note_config

type Config struct {

	// Encryption
	EncryptionKey string `yaml:"encryption_key"`
	UseEncryption bool   `yaml:"use_encryption"`

	// Defaults

	Template  string   `yaml:"template,omitempty"` // Path to template
	Authors   []string `yaml:"authors,omitempty"`  // List of default authors
	BuildMode string   `yaml:"build_mode"`         // Default build mode

	// Build

	BuildCommand []string `yaml:"build_command,omitempty"` // The command that gets run to export the markdown. Replacers {note_path} {output_path} should be used. Replacer {build_file_<nr>} can be used to have access to build files in the build command.
	BuildFiles   []string `yaml:"build_files,omitempty"`   // Names of the files used in the build command, no path just filename.

	BuildFileName    string `yaml:"build_file_name,omitempty"`    // The name of the exported file. Replacers {authors} {subject} {title} can be used.
	BuildAuthor      string `yaml:"build_author,omitempty"`       // How a single author should look in the file name.
	BuildAuthorSplit string `yaml:"build_author_split,omitempty"` // How the authors should be chained together if multiple are present. Replacers {last_name} {middle_names} {fist_name} {given_name} can be used.

	AutoBuild bool `yaml:"auto_build,omitempty"` // Whenever editing or creating a note build it right after.
}
