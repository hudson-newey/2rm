package cli

import (
	"fmt"
	"hudson-newey/2rm/src/models"
	"os"
	"slices"

	// I don't use the flags package provided by the GoLang standard library
	// because I am trying to keep parity with the GNU "rm" command.
	// This is hard when command line parsing doesn't work the same.
	// Therefore, I use a third party package to replicate GNU argument parsing.
	flag "github.com/spf13/pflag"
)

func ParseCliFlags(args []string) models.CliOptions {
	hardDelete := flag.BoolP(HARD_DELETE_CLA, HARD_DELETE_SHORT_CLA, false, "Do not soft-delete FILE(s)")
	softDelete := flag.BoolP(SOFT_DELETE_CLA, SOFT_DELETE_SHORT_CLA, true, "Soft delete a file and store a backup (default /tmp/2rm)")
	bypassProtected := flag.Bool(BYPASS_PROTECTED_CLA, false, "Using this flag will allow you to delete a file protected by the 2rm config")
	overwrite := flag.Bool(OVERWRITE_CLA, false, "Overwrite the disk location with zeros")
	silent := flag.Bool(SILENT_CLA, false, "Do not print out additional information produced by 2rm")
	dryRun := flag.Bool(DRY_RUN_CLA, false, "Perform a dry run and show all files that would be deleted without the dry-run flag")
	shouldNotify := flag.Bool(NOTIFICATION_CLA, false, "Send a system notification once deletion is complete")

	requestingHelp := flag.Bool(HELP_CLA, false, "Display this help and (without deleting anything)")
	requestingVersion := flag.Bool(VERSION_CLA, false, "Output version information (without deleting anything)")

	onlyEmptyDirsDesc := ""
	onlyEmptyDirsShort := flag.BoolP(DIR_CLA, DIR_CLA, false, onlyEmptyDirsDesc)
	onlyEmptyDirsLong := flag.Bool(DIR_CLA_LONG, false, "")
	onlyEmptyDirs := *onlyEmptyDirsShort || *onlyEmptyDirsLong

	verboseDesc := ""
	verboseShort := flag.BoolP(VERBOSE_SHORT_CLA, VERBOSE_SHORT_CLA, false, verboseDesc)
	verboseLong := flag.Bool(VERBOSE_CLA, false, verboseDesc)
	hasVerboseCla := *verboseShort || *verboseLong

	// even though the recursive flag has no effect in 2rm, I add it to the
	// programs list of flags so that 2rm can correctly parse it
	// see: https://github.com/hudson-newey/2rm/issues/27
	flag.BoolP(RECURSIVE_CLA, RECURSIVE_CLA, false, "No action")

	oneFileSystem := flag.Bool(ONE_FILE_SYSTEM_CLA, false, "When a hierarchical structure is being processed, skip any directory that is on a different file system from the one of the corresponding command line argument.")

	// Calling "interactiveWhenParser" calls the flags.Parse() function
	isInteractive, isGroupInteractive, isOnceInteractive := interactiveWhenParser(args)

	return models.CliOptions{
		HardDelete:         *hardDelete,
		SoftDelete:         *softDelete,
		BypassProtected:    *bypassProtected,
		Overwrite:          *overwrite,
		Silent:             *silent,
		DryRun:             *dryRun,
		ShouldNotify:       *shouldNotify,
		RequestingHelp:     *requestingHelp,
		RequestingVersion:  *requestingVersion,
		IsInteractive:      isInteractive,
		IsOnceInteractive:  isOnceInteractive,
		IsGroupInteractive: isGroupInteractive,
		OnlyEmptyDirs:      onlyEmptyDirs,
		Verbose:            hasVerboseCla,
		OneFileSystem:      *oneFileSystem,

		RawArguments: args,
	}
}

// The --interactive flag is annoying because it can act as a boolean flag or
// an "enum like" value.
// If used as a boolean flag, it acts like the standard "-i" flag, where each
// file is prompted for deletion.
//
// However, it is possible to supply a value in the format
// --interactive=[WHEN]
// Where [WHEN] accepts the values
//   - never, no, none
//   - once
//   - always, yes
//
// This is especially annoying because the explicitly provided interactive
// option can conflict with other flags supplied.
//
// e.g. $ 2rm -i --interactive=never
//
// Should deletion be interactive or not? the -i flag explicitly states that
// deletion should be interactive, but he --interactive=never flag contradicts
// the -i flag.
// The GNU "rm" implementation fixes this with a simple solution. Accept the
// last parameter.
// In the example above, the deletion will not be interactive because the -i
// flag is overwritten by the --interactive=never statement.
// The same goes the other way.
//
// e.g. $ 2rm --interactive=never -i
//
// Should run an interactive deletion session.
// This same logic is still applied when multiple --interactive=[WHEN] values
// are provided.
//
// e.g. $ 2rm --interactive=never --interactive=always
//
// Will run an interactive deletion session.
func interactiveWhenParser(args []string) (bool, bool, bool) {
	isInteractive := flag.BoolP(INTERACTIVE_CLA, INTERACTIVE_CLA, false, "")
	isGroupInteractive := flag.BoolP(INTERACTIVE_GROUP_CLA, INTERACTIVE_GROUP_CLA, false, "")

	interactiveWhen := flag.String(INTERACTIVE_WHEN_CLA_PREFIX, "", "")
	flag.Lookup("interactive").NoOptDefVal = "always"
	flag.Parse()

	if *interactiveWhen == "" {
		// if the --interactive flag is used, I set the "isInteractive" to true to
		// pretend like the user has used the "-i" flag
		return *isInteractive, *isGroupInteractive, false
	}

	falseOptions := []string{"never", "no", "none"}
	onceOptions := []string{"once"}
	alwaysOptions := []string{"always", "yes"}

	if slices.Contains(falseOptions, *interactiveWhen) {
		// If the value is "falsy" e.g. --interactive=never, I hack my way through
		// by acting as if the -i flag and the -I flags were never passed.
		// I also return a pointer to a "false" value so that "interactiveOnce" is
		// not set.
		return false, false, false
	}

	if slices.Contains(onceOptions, *interactiveWhen) {
		// Because --interactive=once has overwritten both the -i and -I flags
		// (due to being later in the arg list), I set both the values to falsy.
		//
		// Note: this is not very DRY code, but I think it's a bit more explicit as
		// to what I am doing, and it makes more sense to my current mind. Might
		// change later if I can read code better in the future.
		return false, false, true
	}

	if slices.Contains(alwaysOptions, *interactiveWhen) {
		return true, false, false
	}

	fmt.Println("Invalid argument '", *interactiveWhen, "' for '--", INTERACTIVE_WHEN_CLA_PREFIX, "'")
	fmt.Println("Valid arguments are:")
	fmt.Println("\t- 'never', 'no', 'none'")
	fmt.Println("\t- 'once'")
	fmt.Println("\t- 'always', 'yes'")
	fmt.Println("Try '2rm --help' for more information.")

	os.Exit(1)

	// The go lsp complains if I don't return a value here, even though the
	// program has quit above.
	// If (for some reason) the program gets to this point, I want to just
	// ignore the --interactive[=WHEN] flag
	return false, false, false
}
