package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tiemingo/smn/util"
	"gopkg.in/yaml.v3"
)

type header struct {
	Title   string   `yaml:"title"`
	Authors []string `yaml:"author"`
	Subject string   `yaml:"subtitle"`
}

func (note *Note) BuildNote(buildMode string) (string, error) {

	fileName, err := note.getBuildFileName(note.GetNotePath())
	if err != nil {
		return "", fmt.Errorf("failed to create file name: %v", err)
	}

	// Look for build files
	if buildMode == "" {
		buildMode = note.config.BuildMode
	}
	buildFileReplacers, err := note.getBuildFiles(buildMode)
	if err != nil {
		return "", fmt.Errorf("failed to get build file replacers: %v", err)
	}
	outDir := filepath.Join(note.getNoteDir(), "out")
	buildFileReplacers = append(buildFileReplacers, "{note_path}", note.GetNotePath(), "{output_path}", filepath.Join(outDir, fileName))

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create outDir(%v): %v", outDir, err)
	}

	outputFileName := ""

	// Run build command
	replaceCommand := strings.NewReplacer(buildFileReplacers...)
	buildCommand := note.config.BuildCommand
	for i, commandElement := range buildCommand {

		replacedString, err := util.ReplaceWithHomeDir(replaceCommand.Replace(commandElement))
		if err != nil {
			return "", fmt.Errorf("failed to replace home directory: %v", err)
		}
		buildCommand = slices.Replace(buildCommand, i, i+1, replacedString)
		if strings.Contains(commandElement, "{output_path}") {
			outputFileName = replacedString
		}
	}

	if _, err, stderr := util.RunCommand(buildCommand[0], buildCommand[1:]...); err != nil {
		return "", fmt.Errorf("failed to run %v: %v, stderr: %v", buildCommand, err, stderr)
	}
	return outputFileName, nil
}

// buildFileName returns the filename that should be used for exported notes with all replacer being replaced with actual values
func (note *Note) getBuildFileName(notePath string) (string, error) {

	noteHeader, err := getHeader(notePath)
	if err != nil {
		return "", fmt.Errorf("failed to get header: %v", err)
	}

	// Parse authors
	authors := []string{}
	for _, author := range noteHeader.Authors {
		replacedAuthor := util.ParseAuthor(author)
		authorReplacer := strings.NewReplacer("{last_name}", replacedAuthor.LastName, "{first_name}", replacedAuthor.FirstName, "{given_name}", replacedAuthor.GivenName)
		authors = append(authors, authorReplacer.Replace(note.config.BuildAuthor))
	}
	authorsString := strings.Join(authors, note.config.BuildAuthorSplit)

	titleReplacer := strings.NewReplacer("{authors}", authorsString, "{title}", noteHeader.Title, "{subject}", noteHeader.Subject)
	replacedTitle := titleReplacer.Replace(note.config.BuildFileName)

	// Replace spaces and trim
	replacedTitle = strings.TrimSpace(replacedTitle)
	replacedTitle = strings.ReplaceAll(replacedTitle, " ", note.config.BuildReplaceSpace)

	return replacedTitle, nil
}

// getHeader returns the information from the markdown files yaml header
func getHeader(path string) (header, error) {

	// Load note and extract header
	noteBytes, err := os.ReadFile(path)
	if err != nil {
		return header{}, fmt.Errorf("failed to load note: %v", err)
	}
	parts := strings.SplitN(string(noteBytes), "---", 3)
	if len(parts) < 3 {
		return header{}, fmt.Errorf("failed to find header: %v", err)
	}
	yamlBlock := parts[1]

	// Unmarshal
	var retHeader header
	err = yaml.Unmarshal([]byte(yamlBlock), &retHeader)
	if err != nil {
		return header{}, fmt.Errorf("failed to unmarshal header: %v", err)
	}

	return retHeader, nil
}

func (note *Note) getBuildFiles(buildMode string) ([]string, error) {
	if len(note.config.BuildFiles) == 0 {
		return []string{}, nil
	}

	buildFilePaths := []string{}
	switch buildMode {
	case "n":
		var ok bool
		buildFilePaths, ok = note.getBuildFilePath(note.getNoteDir())
		if !ok {
			return buildFilePaths, fmt.Errorf("failed to find all build files")
		}
	case "s":
		var ok bool
		buildFilePaths, ok = note.getBuildFilePath(note.getNoteDirParent())
		if !ok {
			return buildFilePaths, fmt.Errorf("failed to find all build files")
		}
	case "r":
		currentPath := note.getNoteDir()
		for {
			var ok bool
			buildFilePaths, ok = note.getBuildFilePath(currentPath)
			if ok {
				break
			}

			currentPath = filepath.Dir(currentPath)
			if currentPath == note.getNotesDir() {
				return buildFilePaths, fmt.Errorf("no build files in any parent directory")
			}
		}
	default:
		return buildFilePaths, fmt.Errorf("invalid build mode")
	}
	buildReplacers := []string{}
	for i, path := range buildFilePaths {
		buildReplacers = append(buildReplacers, fmt.Sprintf("{build_file_%v}", i+1), path)
	}
	return buildReplacers, nil
}

func (note *Note) getBuildFilePath(path string) ([]string, bool) {

	files := []string{}
	for _, buidlFile := range note.config.BuildFiles {
		currentPath := filepath.Join(path, buidlFile)
		if stat, err := os.Stat(currentPath); err == nil && !stat.IsDir() {
			files = append(files, currentPath)
		} else {
			return []string{}, false
		}
	}
	return files, true
}
