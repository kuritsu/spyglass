package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetOptions according to arguments
type TargetOptions struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	actions map[string]Command
}

// GetFlags for the current command.
func (o *TargetOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetOptions) Description() string {
	return "Manages targets."
}

// TargetFlags obtains flags for apply action
func TargetFlags() *TargetOptions {
	result := TargetOptions{}
	result.flagSet = flag.NewFlagSet("target", flag.ContinueOnError)
	result.actions = make(map[string]Command)
	result.actions["list"] = TargetListActionFlags(result.flagSet)
	result.actions["update-status"] = TargetUpdateStatusActionFlags(result.flagSet)
	result.flagSet.Usage = func() {
		args := result.flagSet.Args()
		if len(args) > 0 {
			action, ok := result.actions[args[0]]
			if ok {
				action.GetFlags().Usage()
				return
			}
		}
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] <action> ")
		fmt.Println("\nActions:")
		fmt.Println("  list: Paginated list of existing targets.")
		fmt.Println("  update-status: Updates the status of a target.")
		fmt.Println("\nGlobal Flags:")
		result.flagSet.PrintDefaults()
	}
	return &result
}

// Apply the command.
func (o *TargetOptions) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Doing target...")
	nonFlag := o.flagSet.Args()
	if len(nonFlag) == 0 {
		return &runner.ExitError{FlagSet: o.flagSet,
			Error:  errors.New("An action is required"),
			Logger: c.Log,
		}
	}
	switch nonFlag[0] {
	case "list":
		return o.ListTargets()
	case "update-status":
		return o.UpdateStatus()
	}
	return &runner.ExitError{FlagSet: o.flagSet,
		Error:  errors.New("Action not supported"),
		Logger: c.Log,
	}
}

// ListTargets shows a list either formatted or using JSON.
func (o *TargetOptions) ListTargets() runner.Runner {
	return nil
}

// UpdateStatus ensures that the status of the specified target is updated.
func (o *TargetOptions) UpdateStatus() runner.Runner {
	return nil
}
