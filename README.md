# Simple Markdown Notes
Simple markdown notes is a as the name suggests simple markdown notes manager written in go. The main use case of this program is to create notes from templates, edit them and export them to whatever you want, all while your notes get synced using git. I created this program to manage my markdown notes for uni and export them easily as pdfs that I can submit.

## Features
* Automatically synchronize notes using git
* Edit notes in your favorite editor
* Set build command to export the markdown files
* Create notes with templates

## Installation
The program can be installed via go install. 

```bash
go install github.com/tiemingo/smn@latest
```
If you have not installed a golang program before, you need to add the Go bin directory to your systems path.<br>
This usually is your `.bashrc`.
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage
**Create new note**<br>
Note name hast to be set, this will be used as the title. The note can be put into a subject(subfolder) that can be used in the build command. Multiple subjects can be used separated by `/`, this can be helpful for organizing your notes. Please note that only the last subject can be used in the build command. After creating a note it will be opened in the default editor.
```bash
smn create "<subject>/<sub_subject>/<my_note>"
```
**Edit existing note**<br>
A note can be edit by entering the name of the note used to create it. The note will then be opened in the default editor.
```bash
smn edit "<subject>/<sub_subject>/<my_note>"
```
**Remove note**<br>
A note can be deleted by entering the name of the note used to create it.
```bash
smn remove "<subject>/<sub_subject>/<my_note>"
```
**Build note**<br>
A note can be exported by entering the name of the note used to create it. A custom build command can be set, I personally use pandoc and have only tested it with that.
```bash
smn build "<subject>/<sub_subject>/<my_note>"
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

**Get default config**<br>
This command returns a default config that can be pasted to your config file and modified as needed.
```bash
smn config
```

## Synchronization
Git is used for synchronization. The synchronization can be triggered by running `smn sync` or it gets triggered automatically if `auto_sync` is set to true in the config.<br>
The automatic sync triggers before running commands `create`, `edit`, `remove`, `build`, `latest` and after running commands `create`, `edit`.

## Configuration
The editor that opens when a new note is created or an existing note is edited can be set using the `EDITOR` environment variable. The config file itself is pretty self explanatory.<br>
This is the configuration I currently use myself. The config can be set at `<users_config_dir>/smn/config.json`. 
```json
{
	"notes_dir": "~/Documents/notes",
	"template": "~/Documents/notes/templates/template.md",
	"output_dir": "~/Downloads",
	"auto_sync": true,
	"fail_on_sync_error": false,
	"default_authors": ["Tiemingo"],
	"build_command": ["pandoc",  "{note_path}", "~/Documents/notes/build/style.yaml", "-d", "~/Documents/notes/build/proposals.yaml", "-o", "{output_path}.pdf"],
	"build_file_name": "{authors}_{title}",
	"build_author": "{last_name}_{first_name}",
	"build_author_split": "_and_"
}
```
### Replacer
You can use replacer for the build command and the exported files name.
* **Build command:**
  * `{note_path}`: Path of the file that should be exported.
  * `{output_path}`: Output file name without file extension.
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
.
├── .git/
├── build/
│   ├── proposals.yaml
│   └── style.yaml
└── templates/
    └── template.md
```
