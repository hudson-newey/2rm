package cli

import "fmt"

const VERSION = "2rm 0.1.1"

const HELP_DOCS = `
Usage: rm [OPTION]... [FILE]...
Mark FILE(s) for deletion.

"GNU Like" OPTION(s):
-i					Prompt before every deletion request

-r, -R, --recursive			Remove directories and their contents recursively

-I 					Prompt once before deleting more than the interactive
					threshold (default 3)

--interactive[=WHEN]
								never false: never, no, none
								prompt once: once
								always prompt: always, yes

-f, --force				Bypass protections

-v, --verbose				Add additional information to the output

-d, --dir				Only remove empty directories

--help					Display this help and (without deleting anything)

--version				Output version information (without deleting anything)

2rm OPTION(s):

--overwrite				Overwrite the disk location location with zeros

--hard					Do not soft-delete FILE(s)

--soft 					Soft delete a file and a store backup (default /tmp/2rm)

--silent				Do not print out additional information produced by 2rm.
					This is useful for scripting situations

--dry-run				Perform a dry run and show all the files that would be
					deleted if run without the dry-run flag

--bypass-protected			Using this flag will allow you to delete a file
					protected by the 2rm config

--notify 				Send a system notification once deletion is complete

By default, 2rm will soft-delete a file and store a backup (default /tmp/2rm)

Config-based deletion:

You can create a system wide 2rm config file under /etc/2rm/config.yml that will
apply to all users on the system.

Alternatively you can create a config bound to the local user at
~/.local/share/2rm/config.yml.

User config files will always take precedence over system wide configurations.

Using a config will allow you to protect files/directories, set a custom backup
location, and configure directories that should always be hard-deleted.

Exit codes:

1: GNU compatible error
2: An error occured during a 2rm function

Use "man 2rm" for more information.
`

func PrintVersion() {
	fmt.Println(VERSION)
}

func PrintHelp() {
	fmt.Print(HELP_DOCS)
}
