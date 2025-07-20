package models

// TODO: Consolidate with the config model
type CliOptions struct {
	HardDelete         bool
	SoftDelete         bool
	BypassProtected    bool
	Overwrite          bool
	Silent             bool
	DryRun             bool
	ShouldNotify       bool
	RequestingHelp     bool
	RequestingVersion  bool
	IsInteractive      bool // prompt before deleting every file
	IsOnceInteractive  bool // prompt when deleting first file
	IsGroupInteractive bool // prompt when deleting more than 3 files
	OnlyEmptyDirs      bool
	Verbose            bool
	OneFileSystem      bool

	RawArguments []string
}
