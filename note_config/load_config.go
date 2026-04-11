package note_config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

var defaultConfig = Config{
	UseEncryption: true,

	Template:     "~/Documents/notes-template.tmpl",
	Authors:      []string{"Tiemingo"},
	BuildMode:    "r",
	GitIgnoreOut: true,

	BuildCommand: []string{
		"pandoc",
		"{note_path}",
		"{build_file_1}",
		"-d",
		"{build_file_2}",
		"-o",
		"{output_path}.pdf",
	},
	BuildFiles: []string{
		"style.yaml",
		"proposals.yaml",
	},
	BuildFileName:     "{authors} {title}",
	BuildAuthor:       "{last_name} {first_name}",
	BuildAuthorSplit:  " and ",
	BuildReplaceSpace: "_",

	AutoBuild: true,
}

// GetDefaultConfig returns the default config.
func GetDefaultConfig() string {
	byteConfig, err := yaml.Marshal(defaultConfig)
	if err != nil {
		log.Printf("WARN: Failed to marshal default config: %v\n", err)
	}
	return string(byteConfig)
}

func GetConfig(baseDir, relDir string) (Config, error) {

	// Create list of paths that could have a config file
	configPaths := []string{}
	relPaths := strings.Split(strings.Trim(filepath.ToSlash(relDir), "/"), "/")
	for i := range relPaths {
		fullPath := slices.Concat([]string{baseDir}, relPaths[:i+1], []string{"config.yaml"})
		configPaths = append(configPaths, filepath.Join(fullPath...))
	}

	var config Config

	if len(configPaths) == 0 {
		return config, fmt.Errorf("failed to get any config paths")
	}

	for i, configPath := range configPaths {

		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return config, fmt.Errorf("failed to read config file: %v", err)
			} else if errors.Is(err, os.ErrNotExist) && i == 0 {
				return config, fmt.Errorf("there has to be a config file in the base note dir: %v", err)
			}
			continue
		}

		// Overwrite existing config with options of higher priority
		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			return config, fmt.Errorf("failed to unmarshal config file: %v", err)
		}
	}

	return config, nil
}
