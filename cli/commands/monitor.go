package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// MonitorOptions according to arguments
type MonitorOptions struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	actions map[string]Command
}

// GetFlags for the current command.
func (o *MonitorOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *MonitorOptions) Description() string {
	return "Manages monitors."
}

// MonitorFlags obtains flags for apply action
func MonitorFlags() *TargetOptions {
	result := TargetOptions{}
	result.flagSet = flag.NewFlagSet("monitor", flag.ContinueOnError)
	result.actions = make(map[string]Command)
	result.actions["add"] = MonitorAddActionFlags(result.flagSet)
	// result.actions["get"] = TargetGetActionFlags(result.flagSet)
	result.actions["list"] = MonitorListActionFlags(result.flagSet)
	// result.actions["rm"] = TargetRemoveActionFlags(result.flagSet)
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
		fmt.Println("  spyglass monitor [global-flags] <action> ")
		fmt.Println("\nActions:")
		for _, k := range GetSortedKeyList(result.actions) {
			fmt.Printf("  %s: %s\n", k, result.actions[k].Description())
		}
		fmt.Println("\nGlobal Flags:")
		result.flagSet.PrintDefaults()
	}
	return &result
}

// Apply the command.
func (o *MonitorOptions) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Doing monitor...")
	nonFlag := o.flagSet.Args()
	if len(nonFlag) == 0 {
		return &runner.ExitError{FlagSet: o.flagSet,
			Error:  errors.New("An action is required"),
			Logger: c.Log,
		}
	}
	subcommand, ok := o.actions[nonFlag[0]]
	if !ok {
		return &runner.ExitError{FlagSet: o.flagSet,
			Error:  errors.New("Action not supported"),
			Logger: c.Log,
		}
	}
	o.actions[nonFlag[0]].GetFlags().Parse(nonFlag[1:])
	return subcommand.Apply(c)
}
