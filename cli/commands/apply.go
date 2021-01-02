package commands

import (
	"flag"
)

// ApplyOptions according to arguments
type ApplyOptions struct {
	flagSet   *flag.FlagSet
	Recursive bool
}

// Apply the configuration in the given directory.
func (o *ApplyOptions) Apply(c *CommandLineContext) {
	c.Log.Debug("Executing apply.")
}

// GetFlags for the current command.
func (o *ApplyOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *ApplyOptions) Description() string {
	return "Apply the configuration contained in *.sgc files in the specified directory."
}

// ApplyFlags obtains flags for apply action
func ApplyFlags() *ApplyOptions {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	result := ApplyOptions{flagSet: fs}
	fs.BoolVar(&result.Recursive, "r", false, "Scan specified path recursively for config files.")
	return &result
}
