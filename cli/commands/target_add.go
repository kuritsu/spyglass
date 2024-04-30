package commands

import (
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetAddAction represents the list action of targets.
type TargetAddAction struct {
	flagSet  *flag.FlagSet
	c        *CommandLineContext
	fileName string
}

// GetFlags for the current command.
func (o *TargetAddAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetAddAction) Description() string {
	return "Adds a new target with its children."
}

// TargetAddActionFlags obtains flags for target list action
func TargetAddActionFlags(parentFs *flag.FlagSet) *TargetAddAction {
	result := TargetAddAction{}
	result.flagSet = flag.NewFlagSet("add", flag.ContinueOnError)
	result.flagSet.StringVar(&result.fileName, "f", "", "File name (must end in json, yaml or yml) of the target data.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] add -f fileName.json")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *TargetAddAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply target add...")
	target, err := types.NewTargetFromFile(o.fileName)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	err = c.Caller.InsertOrUpdateTarget(target, true)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	o.writeTarget(target)
	return nil
}

func (o *TargetAddAction) writeTarget(target types.TargetRef) {
	o.c.Log.Info("Created target ", target.ID)
	if target.Children == nil {
		return
	}
	for _, t := range target.Children {
		o.writeTarget(t)
	}
}
