package commands

import (
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/runner"
)

// RoleAddAction represents the list action of targets.
type RoleAddAction struct {
	flagSet  *flag.FlagSet
	c        *CommandLineContext
	fileName string
}

// GetFlags for the current command.
func (o *RoleAddAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *RoleAddAction) Description() string {
	return "Adds a new role."
}

// RoleAddActionFlags obtains flags for target list action
func RoleAddActionFlags(parentFs *flag.FlagSet) *RoleAddAction {
	result := RoleAddAction{}
	result.flagSet = flag.NewFlagSet("add", flag.ContinueOnError)
	result.flagSet.StringVar(&result.fileName, "f", "", "File name (must end in json, yaml or yml) of the target data.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass role [global-flags] add -f fileName.json")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *RoleAddAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply role add...")
	role, err := api.NewObjectFromFile[types.Role](o.fileName)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	err = c.Caller.InsertRole(role)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	return nil
}
