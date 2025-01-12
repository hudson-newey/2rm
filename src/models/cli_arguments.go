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
	IsInteractive      bool
	IsGroupInteractive bool
	OnlyEmptyDirs      bool
	Verbose            bool

	RawArguments []string
}
