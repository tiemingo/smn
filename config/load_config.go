package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var defaultConfig = Config{
	NotesDir:   "~/Documents/notes",
	BuildFiles: "~/Documents/notes/build",
	Template:   "", // No template
	OutputDir:  "~/Downloads",

	AutoSync:        true,
	FailOnSyncError: false,

	DefaultAuthors: "[Tiemingo]",

	BuildCommand:     "pandoc %v build/style.yaml -d build/proposals.yaml -o %v.pdf",
	BuildFileName:    "{authors}_{title}",
	BuildAuthorSplit: "_and_",
	BuildAuthor:      "{last_name}_{fist_name}",
}

// GetDefaultConfig returns the default config.
func GetDefaultConfig() string {
	byteConfig, err := json.MarshalIndent(defaultConfig, "", "	")
	if err != nil {
		log.Printf("WARN: Failed to marshal default config: %v\n", err)
	}
	return string(byteConfig)
}

// GetConfig returns the loaded config or if it can't be found or parsed, the default config.
func GetConfig() Config {

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("WARN: Failed to get config directory: %v\n", err)
		return defaultConfig
	}

	// Load config file
	file, err := os.Open(filepath.Join(configDir, "smn", "config.json"))
	if err != nil {
		log.Printf("WARN: Failed to load config: %v\n", err)
		return defaultConfig
	}
	defer file.Close()

	// Decode json and parse
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Printf("WARN: Failed to parse config: %v\n", err)
		return defaultConfig
	}

	return config
}
