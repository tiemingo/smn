# Simple Markdown Notes
Simple markdown notes is a as the name suggests simple markdown notes manager written in go. The main use case of this program is to create notes from templates, edit them and export them to whatever you want, all while your notes get synced using git. I created this program to manage my markdown notes for uni and export them easily as pdfs that I can submit.

## Features
* Automatically synchronize notes using git
* Support for multiple git repos
* Edit notes in your favorite editor
* Set build command to export the markdown files
* Create notes with templates
* Automatically reexport markdown file upon editing
* Modular configuration for every topic, subject and note

## Installation
The program can be installed via go install. 

```bash
go install -tags release github.com/tiemingo/smn@latest
```
If you have not installed a golang program before, you need to add the Go bin directory to your systems path.<br>
This usually is your `.bashrc`.
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Setup
You should first create a notes directory where ever you want, note that this directory should not contain a .git folder.<br>
The direct child directories of the notes directory are called topic, they should be a git repo if you want them to be synced. They each have to contain a `config.yaml`. Check the [suggested structure](#suggested-notes-directory-structure) down below.
Next you have to create a `config.json` in your user config folder for `smn`. In this config file you set the path to your newly created notes directory, an encryption key for AES-256 if you want to use encryption for some notes and the syncing behavior. See [configuration](#configuration) for the exact fields you can set in the config. 


## Usage
**Create new note**<br>
A note has to have at least a topic and a name. Subjects are subdirectories the note will be stored in. Topic, each of the subjects and name have to be separated by `/`. Subjects can be really helpful for organizing your notes. Please note that only the last subject can be used in the build command. After creating a note it will be opened in the default editor.
```bash
smn create "<topic>/<subject>/<sub_subject>/<my_note>"
```
**Edit existing note**<br>
A note can be edit by entering the name of the note used to create it. The note will then be opened in the default editor.
```bash
smn edit "<topic>/<subject>/<sub_subject>/<my_note>"
```
**Remove note**<br>
A note can be deleted by entering the name of the note used to create it.
```bash
smn remove "<topic>/<subject>/<sub_subject>/<my_note>"
```
**Build note**<br>
A note can be exported by entering the name of the note used to create it. A custom build command can be set, I personally use pandoc and have only tested it with that. <br>
Optionally you can overwrite the config files build mode, by adding `-b` and specifying the build mode behind it. The build command returns the filepath to the last outfile if multiple are specified.
```bash
smn build "<topic>/<subject>/<sub_subject>/<my_note>" -b <build_mode>
```
**List latest notes**<br>
The latest modified notes can be listed. The amount of the returned notes can be set, if the amount is <= 0 all notes sorted by last modified are returned. If the amount provided is higher than the actual amount of notes, all notes are returned.
```bash
smn latest <amount>
```

**Sync notes**<br>
This command checks for new commits in the remote, merges local commits and pushes the changes to the remote.
```bash
smn sync
```

**Get default config** <br>
This command returns a default config for the notes and a general config for smn that can be pasted to your config files and modified as needed.
```bash
smn config
```

## Synchronization
Git is used for synchronization. The synchronization can be triggered by running `smn sync` or it gets triggered automatically if `auto_sync` is set to true in the config.<br>
The automatic sync triggers before running commands `create`, `edit`, `remove`, `build`, `latest` and after running commands `create`, `edit`.

## Exporting
You can either export notes manually by running the `build` command or set `auto_build` in your `config.yaml` to true. This will build the markdown every time a note is created or edited.
### Build files
Build files can be set in the `config.yaml`. The build files should only be the filenames, no path. The path later gets resolved. The method of resolving can be set using build mode flag in command or field in config.
### Build mode
There are 3 build mode:
 - `n` looks for build files in the note's directory.
 - `s` looks for build files in the note subject directory.
 - `r` iterates through the parent directories until it finds build files, starting in the note's directory.

## Configuration
**General configuration of smn:** <br>
The editor that opens when a new note is created or an existing note is edited can be set using the `EDITOR` environment variable.<br>
For the configuration a `config.json` is used that can be set at `<users_config_dir>/smn/config.json`.<br>
This is an example configuration file with all options.
```json
{
	"encryption_key": "AES-256 encryption key",
	"notes_dir": "~/Documents/notes",
	"auto_sync": true,
	"fail_on_sync_error": false
}
```

**Configuration of topics, subjects and notes:** <br>
The configurations for a note are applied in order from topic, to subjects to note. This means if you set a option in your topic `config.yaml` and set a different value in the `config.yaml` in the directory of your note, the value directly in the note's directory takes priority. If a value is not present in a more higher priority config, the one from a the lower priority config is used.<br>
This is really useful, if you want to change for example the formatting of the output filename for a certain subject, but want to use the topics formatting option in all other subjects. Then you only have to set one value in the `config.yaml` in the subject directory.<br>
This is an example configuration file with all options.
```yaml
use_encryption: true

template: ~/Documents/notes-template.tmpl # Default body when creating a note
authors:            # Authors set in the header when creating a new note 
    - Tiemingo      
build_mode: r       # Build mode used for selection of build files
gitignore_out: true # This adds the /out directory of a newly created note to the notes .gitignore

build_command: # Command used to export markdown
    - pandoc
    - '{note_path}'
    - '{build_file_1}'
    - -d
    - '{build_file_2}'
    - -o
    - '{output_path}.pdf'
build_files: # Files that can be used on build command, only name (path gets resolved using build mode)
    - style.yaml
    - proposals.yaml
build_file_name: '{authors} {title}' # Exported files filename
build_author: '{last_name} {first_name}'  # Formatting of each author
build_author_split: ' and ' # Separator between authors
build_replace_space: "_" # String that replaces spaces in the out filename. Use " " to leave spaces. 
auto_build: true # Should automatically export note on change
```

### Replacer
You can use replacer for the build command and the exported files name.
* **Build command:**
  * `{note_path}`: Path of the file that should be exported.
  * `{output_path}`: Output file name without file extension.
  * `{build_file_<nr>}`: Resolved path to the build file specified in config. With nr being the stop in the list starting with 1.
* **Build file name:**
  * `{authors}`: Authors set in the markdown header.
  * `{subject}`: Subtitle in the markdown header. By default this is the directory the markdown file is in.
  * `{title}`: Title set in the markdown header.
* **Author:**
  * `{first_name}`: First name of an author.
  * `{given_name}`: Given name of an author.
  * `{last_name}`: Last name of an author.


## Suggested notes directory structure
```text
notes/
├── uni (this is a topic and should be a git repo if you want to sync it)
│   ├── config.yaml (config file that is used for all note in the uni topic)
│   ├── Oop (a subject)
│   ├── Numerik (a subject)
│   │   ├── config.yaml (the config options here replace the ones set in the uni config)
│   │   └── style.yaml  (this is a differnt styling that can be used in the Numerik subject, depending on build mode)
│   ├── proposals.yaml
│   └── style.yaml
└── private (this is a different topic, should be a git repo aswell)
    ├── config.yaml
    ├── proposals.yaml
    └── style.yaml
```

## Small tips and inspirations
In this section i want to show you some commands you can use for a better workflow with your notes. Feel free to take this suggestions and use them to expand and simplify your smn workflow.<br>
**Selecting note to edit:** <br>
This command used `dmenu` to list the last 10 notes and opens the selected note to be edited.
```bash
smn edit "$(smn latest | dmenu)"
```
**Selecting note to build and open it right after:** <br>
This command uses `dmenu` to display all your notes, when you select one it gets exported and opened on your default application using `xdg-open`.
```bash
xdg-open "$(smn build "$(smn latest | dmenu)")"
```
As you can see those commands are quite similar, just experiment with them a bit and find some cool use cases. Good luck!
