package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/tiemingo/smn/config"
	"github.com/tiemingo/smn/note_config"
)

func main() {

	parser, createCmd, createPath, editCmd, editPath, removeCmd, removePath, buildCmd, buildPath, buildMode, latestCmd, latestAmount, syncCmd, configCmd := SetupParser()

	// Parse the arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// Route the commands
	if createCmd.Happened() {
		err = createNote(*createPath)
	} else if editCmd.Happened() {
		err = editNote(*editPath)
	} else if removeCmd.Happened() {
		err = removeNote(*removePath)
	} else if buildCmd.Happened() {
		bm := ""
		if buildMode != nil {
			bm = *buildMode
		}
		err = buildNote(*buildPath, bm)
	} else if latestCmd.Happened() {
		err = latestNotes(*latestAmount)
	} else if syncCmd.Happened() {
		err = syncNotes()
	} else if configCmd.Happened() {
		fmt.Printf("config.json in config directory:\n\n%v\n\n\nconfig.yaml in notes directories:\n\n%v\n", config.GetDefaultConfig(), note_config.GetDefaultConfig())

	}

	if err != nil {
		log.Fatalf("failed with err: %v", err)
	}

}

// SetupParser creates all commands available and returns them with their respective arguments
func SetupParser() (*argparse.Parser, *argparse.Command, *string, *argparse.Command, *string, *argparse.Command, *string, *argparse.Command, *string, *string, *argparse.Command, *int, *argparse.Command, *argparse.Command) {
	parser := argparse.NewParser("smn", "A simple markdown note manager")

	// Setup Add
	createCmd := parser.NewCommand("create", "Create a new note")
	addPath := createCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup Edit
	editCmd := parser.NewCommand("edit", "Edit an existing note")
	editPath := editCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup Remove
	removeCmd := parser.NewCommand("remove", "Remove a note")
	removePath := removeCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup build
	buildCmd := parser.NewCommand("build", "Export a note")
	buildPath := buildCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})
	buildMode := buildCmd.String("b", "build-mode", &argparse.Options{Required: false, Help: "Build mode to overwrite default (e.g. s)"})

	// Setup latest
	latestCmd := parser.NewCommand("latest", "Get most recent created or edited notes, use 0 to get all notes")
	latestAmount := latestCmd.IntPositional(&argparse.Options{Required: false, Help: "5"})

	// Setup Sync
	syncCmd := parser.NewCommand("sync", "Sync notes with remote")

	// Setup Config
	configCmd := parser.NewCommand("config", "Configure the app")

	return parser, createCmd, addPath, editCmd, editPath, removeCmd, removePath, buildCmd, buildPath, buildMode, latestCmd, latestAmount, syncCmd, configCmd
}
