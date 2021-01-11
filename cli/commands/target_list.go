package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kuritsu/spyglass/cli/runner"
)

// TargetListAction represents the list action of targets.
type TargetListAction struct {
	flagSet   *flag.FlagSet
	c         *CommandLineContext
	filter    string
	format    string
	pageIndex int
	pageSize  int
}

// GetFlags for the current command.
func (o *TargetListAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *TargetListAction) Description() string {
	return "Paginated list of targets."
}

// TargetListActionFlags obtains flags for target list action
func TargetListActionFlags(parentFs *flag.FlagSet) *TargetListAction {
	result := TargetListAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.filter, "f", "", "Substring the target IDs must contain.")
	result.flagSet.StringVar(&result.format, "o", "json", "Output format. Can be json.")
	result.flagSet.IntVar(&result.pageIndex, "pi", 0, "Page index.")
	result.flagSet.IntVar(&result.pageSize, "ps", 10, "Page size.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass target [global-flags] list [flags]")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *TargetListAction) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Apply target list...")
	result, err := c.Caller.ListTargets(o.filter, o.pageIndex, o.pageSize)
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
