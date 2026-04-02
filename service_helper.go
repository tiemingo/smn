package main

import (
	"log"

	"github.com/tiemingo/smn/config"
)

// syncIfWanted syncs the notes if the config option is set to auto sync.
// If an error occurs while syncing and the config option fail on sync error is enabled, the program exits with an error.
func syncIfWanted(cfg config.Config, optionalCommitMessage ...string) error {

	// Check if sync is wanted
	if !cfg.AutoSync {
		return nil
	}

	err := syncNotes(optionalCommitMessage...)

	// Check if should exit
	if cfg.FailOnSyncError {
		log.Fatalf("failed to sync, if you don't want the program to exit on failed sync, you can change it in the config: %v", err)
	}

	return err
}
