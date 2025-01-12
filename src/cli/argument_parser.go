package cli

import (
	"flag"
	"hudson-newey/2rm/src/models"
)

func ParseCliFlags(args []string) models.CliOptions {
	hardDelete := flag.Bool(HARD_DELETE_CLA, false, "Do not soft-delete FILE(s)")
	softDelete := flag.Bool(SOFT_DELETE_CLA, true, "Soft delete a file and store a backup (default /tmp/2rm)")
	bypassProtected := flag.Bool(BYPASS_PROTECTED_CLA, false, "Using this flag will allow you to delete a file protected by the 2rm config")
	overwrite := flag.Bool(OVERWRITE_CLA, false, "Overwrite the disk location with zeros")
	silent := flag.Bool(SILENT_CLA, false, "Do not print out additional information produced by 2rm")
	dryRun := flag.Bool(DRY_RUN_CLA, false, "Perform a dry run and show all files that would be deleted without the dry-run flag")
	shouldNotify := flag.Bool(NOTIFICATION_CLA, false, "Send a system notification once deletion is complete")

	isInteractive := flag.Bool(INTERACTIVE_CLA, false, "")
	isGroupInteractive := flag.Bool(INTERACTIVE_GROUP_CLA, false, "")

	requestingHelp := flag.Bool(HELP_CLA, false, "Display this help and (without deleting anything)")
	requestingVersion := flag.Bool(VERSION_CLA, false, "Output version information (without deleting anything)")

	onlyEmptyDirsDesc := ""
	onlyEmptyDirsShort := flag.Bool(DIR_CLA, false, onlyEmptyDirsDesc)
	onlyEmptyDirsLong := flag.Bool(DIR_CLA_LONG, false, "")
	onlyEmptyDirs := *onlyEmptyDirsShort || *onlyEmptyDirsLong

	verboseDesc := ""
	verboseShort := flag.Bool(VERBOSE_SHORT_CLA, false, verboseDesc)
	verboseLong := flag.Bool(VERBOSE_CLA, false, "")
	hasVerboseCla := *verboseShort || *verboseLong

	flag.Parse()

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
		IsInteractive:      *isInteractive,
		IsGroupInteractive: *isGroupInteractive,
		OnlyEmptyDirs:      onlyEmptyDirs,
		Verbose:            hasVerboseCla,

		RawArguments: args,
	}
}
