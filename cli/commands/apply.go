package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/sgc"
)

// ApplyOptions according to arguments
type ApplyOptions struct {
	flagSet   *flag.FlagSet
	Recursive bool
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
	fs.BoolVar(&result.Recursive, "r", false, "Scan specified paths recursively for config files.")
	fs.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass apply [flags] path1 [path2] ...")
		fmt.Println("\nFlags:")
		fs.PrintDefaults()
	}
	return &result
}

// Apply the configuration in the given directory.
func (o *ApplyOptions) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Executing apply.")
	paths := o.flagSet.Args()
	if len(paths) == 0 {
		return &runner.ExitError{
			Error:   errors.New("A path is required"),
			FlagSet: o.flagSet,
			Logger:  c.Log,
		}
	}
	fileList := []*sgc.File{}
	for _, d := range paths {
		files, err := c.SgcManager.GetFiles(d, o.Recursive)
		if err != nil {
			c.Log.Error(err)
			return &runner.ExitError{
				Error:  err,
				Logger: c.Log,
			}
		}
		fileList = append(fileList, files...)
	}
	c.Log.Debug("Processing ", len(fileList), " files...")
	return nil
}
