package cli

// 2rm CLI arguments
const HARD_DELETE_CLA = "hard"
const HARD_DELETE_SHORT_CLA = "H"
const SOFT_DELETE_CLA = "soft"
const SOFT_DELETE_SHORT_CLA = "S"
const SILENT_CLA = "silent"
const DRY_RUN_CLA = "dry-run"
const BYPASS_PROTECTED_CLA = "bypass-protected"
const OVERWRITE_CLA = "overwrite"
const NOTIFICATION_CLA = "notify"

const VERBOSE_CLA = "verbose"
const VERBOSE_SHORT_CLA = "v"

// gnu rm CLI arguments
const INTERACTIVE_CLA = "i"
const INTERACTIVE_GROUP_CLA = "I"

const DIR_CLA = "d"
const DIR_CLA_LONG = "dir"

const HELP_CLA = "help"
const VERSION_CLA = "version"

// while this flag has no effect, I have added it as a supported cli argument
// to maintain full backwards compatibility with the GNU rm command
// see: https://github.com/hudson-newey/2rm/issues/27
const RECURSIVE_CLA = "r"
const ONE_FILE_SYSTEM_CLA = "one-file-system"

// This is a special command line argument because it takes a value
// This is why this this cla has _PREFIX in the variable name
const INTERACTIVE_WHEN_CLA_PREFIX = "interactive"

// TODO: Remove this
var SupportedCliArguments = []string{
	HARD_DELETE_CLA,
	HARD_DELETE_SHORT_CLA,
	SOFT_DELETE_CLA,
	SOFT_DELETE_SHORT_CLA,
	SILENT_CLA,
	DRY_RUN_CLA,
	BYPASS_PROTECTED_CLA,
	OVERWRITE_CLA,
	NOTIFICATION_CLA,
	VERBOSE_CLA,
	VERBOSE_SHORT_CLA,
	INTERACTIVE_CLA,
	INTERACTIVE_GROUP_CLA,
	DIR_CLA,
	DIR_CLA_LONG,
	HELP_CLA,
	VERBOSE_CLA,
	RECURSIVE_CLA,
	ONE_FILE_SYSTEM_CLA,

	// TODO: remove this hack for --interactive accepting both boolean and an
	// enum of possible values
	INTERACTIVE_WHEN_CLA_PREFIX,
	INTERACTIVE_WHEN_CLA_PREFIX + "=never",
	INTERACTIVE_WHEN_CLA_PREFIX + "=no",
	INTERACTIVE_WHEN_CLA_PREFIX + "=none",
	INTERACTIVE_WHEN_CLA_PREFIX + "=once",
	INTERACTIVE_WHEN_CLA_PREFIX + "=always",
	INTERACTIVE_WHEN_CLA_PREFIX + "=yes",
}
