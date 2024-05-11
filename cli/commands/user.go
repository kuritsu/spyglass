package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// UserOptions according to arguments
type UserOptions struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	actions map[string]Command
}

// GetFlags for the current command.
func (o *UserOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *UserOptions) Description() string {
	return "Manages users."
}

// UserFlags obtains flags for apply action
func UserFlags() *UserOptions {
	result := UserOptions{}
	result.flagSet = flag.NewFlagSet("user", flag.ContinueOnError)
	result.actions = make(map[string]Command)
	result.actions["list"] = UserListActionFlags(result.flagSet)
	result.actions["token"] = UserTokenActionFlags(result.flagSet)
	//result.actions["rm"] = TargetUpdateStatusActionFlags(result.flagSet)
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
		fmt.Println("  spyglass user [global-flags] <action> ")
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
func (o *UserOptions) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Doing user...")
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
