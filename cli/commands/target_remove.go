package commands

import (
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetRemoveAction represents the remove action for a target and descendants.
type TargetRemoveAction struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	id      string
}

// GetFlags for the current command.
func (o *TargetRemoveAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetRemoveAction) Description() string {
	return "Removes a target and its children."
}

// TargetRemoveActionFlags obtains flags for target remove action
func TargetRemoveActionFlags(parentFs *flag.FlagSet) *TargetRemoveAction {
	result := TargetRemoveAction{}
	result.flagSet = flag.NewFlagSet("rm", flag.ContinueOnError)
	result.flagSet.StringVar(&result.id, "id", "", "Target ID.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] rm -id my_target/child")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *TargetRemoveAction) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Apply target remove...")
	result, err := c.Caller.DeleteTarget(o.id)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	fmt.Printf("%v targets removed.\n", result)
	return nil
}
