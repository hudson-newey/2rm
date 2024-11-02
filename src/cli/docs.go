package cli

import "fmt"

const VERSION = "2rm 0.1.0"

const HELP_DOCS = `
Usage: rm [OPTION]... [FILE]...
Mark FILE(s) for deletion.

"GNU Like" OPTION(s):
-i					Prompt before every deletion request
--force					Bypass 2rm protections

2rm OPTION(s) Flags:
--overwrite				Overwrite the disk location location with zeros
--hard					Do not soft-delete FILE(s)
--soft 					Soft delete a file and a store backup (default /tmp/2rm)
--silent				Do not print out additional information priduced by 2rm. This is useful for scripting situations
--dry-run				Perform a dry run and show all the files that would be deleted
--bypass-protected		Using this flag will allow you to delete a file protected by the 2rm config
--notify 				Send a system notification once deletion is complete

By default, 2rm will soft-delete a file and store a backup (default /tmp/2rm)

You can create a 2rm config file under ~/.local/share/2rm/config.yml that will
allow you to protect files/directories, set a custom backup location, and
configure directories that should always be hard-deleted.

Use "man 2rm" for more information.
`

func PrintVersion() {
	fmt.Println(VERSION)
}

func PrintHelp() {
	fmt.Print(HELP_DOCS)
}