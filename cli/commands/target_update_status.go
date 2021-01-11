package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetUpdateStatusAction represents the list action of targets.
type TargetUpdateStatusAction struct {
	flagSet           *flag.FlagSet
	c                 *CommandLineContext
	id                string
	status            int
	statusDescription string
}

// GetFlags for the current command.
func (o *TargetUpdateStatusAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetUpdateStatusAction) Description() string {
	return "Updates the status of a target."
}

// TargetUpdateStatusActionFlags obtains flags for target list action
func TargetUpdateStatusActionFlags(parentFs *flag.FlagSet) *TargetUpdateStatusAction {
	result := TargetUpdateStatusAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.id, "id", "", "[MANDATORY] Target ID.")
	result.flagSet.IntVar(&result.status, "s", 0, "[MANDATORY] Status (0 - 100).")
	result.flagSet.StringVar(&result.statusDescription, "d", "", "Status description.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] update-status [flags]")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *TargetUpdateStatusAction) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Apply target update-status...")
	idFound := false
	statusFound := false
	o.flagSet.Visit(func(f *flag.Flag) {
		if f.Name == "id" {
			idFound = true
		}
		if f.Name == "s" {
			statusFound = true
		}
	})
	if !idFound {
		return &runner.ExitError{Error: errors.New("id (target ID) flag is required"),
			Logger: c.Log, FlagSet: o.flagSet}
	}
	if !statusFound {
		return &runner.ExitError{Error: errors.New("s (status) flag is required"),
			Logger: c.Log, FlagSet: o.flagSet}
	}
	err := c.Caller.UpdateTargetStatus(o.id, o.status, o.statusDescription)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	c.Log.Debug("Status updated.")
	return nil
}
