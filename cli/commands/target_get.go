package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetGetAction represents the get action for target a given target and its direct descendants.
type TargetGetAction struct {
	flagSet         *flag.FlagSet
	c               *CommandLineContext
	id              string
	format          string
	includeChildren bool
}

// GetFlags for the current command.
func (o *TargetGetAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetGetAction) Description() string {
	return "Paginated list of targets."
}

// TargetGetActionFlags obtains flags for target get action
func TargetGetActionFlags(parentFs *flag.FlagSet) *TargetGetAction {
	result := TargetGetAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.id, "id", "", "Target ID.")
	result.flagSet.StringVar(&result.format, "f", "json", "Output format. Allowed values: json, yaml.")
	result.flagSet.BoolVar(&result.includeChildren, "c", false, "Include children.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] get -id my_target/children -c")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *TargetGetAction) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Apply target get...")
	result, err := c.Caller.GetTargetByID(o.id, o.includeChildren)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	switch o.format {
	case "json":
		c.Log.Debug("Displaying targets as json...")
		jsonBytes, _ := json.Marshal(result)
		var out bytes.Buffer
		json.Indent(&out, jsonBytes, "", "  ")
		out.WriteTo(os.Stdout)
	}
	return nil
}
