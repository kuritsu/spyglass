package commands

import (
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/cli/runner"
)

// MonitorListAction represents the list action of monitors.
type MonitorListAction struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	PaginatedAction
}

// GetFlags for the current command.
func (o *MonitorListAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *MonitorListAction) Description() string {
	return "List monitors."
}

// MonitorListActionFlags obtains flags for target list action
func MonitorListActionFlags(parentFs *flag.FlagSet) *MonitorListAction {
	result := MonitorListAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.format, "o", "json", "Output format. Can be json.")
	result.flagSet.IntVar(&result.pageIndex, "pi", 0, "Page index.")
	result.flagSet.IntVar(&result.pageSize, "ps", 10, "Page size.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass monitor [global-flags] list [flags]")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *MonitorListAction) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Apply target list...")
	result, err := c.Caller.ListMonitors(o.pageIndex, o.pageSize)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	switch o.format {
	case "json":
		c.Log.Debug("Displaying monitors as json...")
		DisplayJson(result)
	}
	return nil
}
