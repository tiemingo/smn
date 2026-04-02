package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/tiemingo/smn/config"
)

func main() {

	parser, addCmd, addPath, editCmd, editPath, buildCmd, buildPath, latestCmd, latestAmount, syncCmd, configCmd := SetupParser()

	// Parse the arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// Route the commands
	if addCmd.Happened() {
		err = addNote(*addPath)
	} else if editCmd.Happened() {
		err = editNote(*editPath)
	} else if buildCmd.Happened() {
		err = buildNote(*buildPath)
	} else if latestCmd.Happened() {
		err = latestNotes(*latestAmount)
	} else if syncCmd.Happened() {
		err = syncNotes()
	} else if configCmd.Happened() {
		fmt.Println(config.GetDefaultConfig())
	}

	if err != nil {
		log.Fatalf("failed with err: %v", err)
	}

}

func SetupParser() (*argparse.Parser, *argparse.Command, *string, *argparse.Command, *string, *argparse.Command, *string, *argparse.Command, *int, *argparse.Command, *argparse.Command) {
	parser := argparse.NewParser("smn", "A simple markdown note manager")

	// Setup Add
	addCmd := parser.NewCommand("add", "Create a new note")
	addNote := addCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup Edit
	editCmd := parser.NewCommand("edit", "Edit an existing note")
	editNote := editCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup build
	buildCmd := parser.NewCommand("build", "Export a note")
	buildNote := buildCmd.StringPositional(&argparse.Options{Required: true, Help: "<subfolder/my_note_title>"})

	// Setup latest
	latestCmd := parser.NewCommand("latest", "Get most recent created or edited notes, use 0 to get all notes")
	latestAmount := latestCmd.IntPositional(&argparse.Options{Required: false, Help: "5"})

	// Setup Sync
	syncCmd := parser.NewCommand("sync", "Sync notes with remote")

	// Setup Config
	configCmd := parser.NewCommand("config", "Configure the app")

	return parser, addCmd, addNote, editCmd, editNote, buildCmd, buildNote, latestCmd, latestAmount, syncCmd, configCmd
}
