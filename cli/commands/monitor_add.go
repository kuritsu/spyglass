package commands

import (
	"flag"
	"fmt"

	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/runner"
)

// MonitorAddAction represents the list action of targets.
type MonitorAddAction struct {
	flagSet  *flag.FlagSet
	c        *CommandLineContext
	fileName string
}

// GetFlags for the current command.
func (o *MonitorAddAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *MonitorAddAction) Description() string {
	return "Adds a new monitor."
}

// MonitorAddActionFlags obtains flags for target list action
func MonitorAddActionFlags(parentFs *flag.FlagSet) *MonitorAddAction {
	result := MonitorAddAction{}
	result.flagSet = flag.NewFlagSet("add", flag.ContinueOnError)
	result.flagSet.StringVar(&result.fileName, "f", "", "File name (must end in json, yaml or yml) of the monitor data.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass monitor [global-flags] add -f fileName.json")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *MonitorAddAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply monitor add...")
	monitor, err := api.NewObjectFromFile[types.Monitor](o.fileName)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	err = c.Caller.InsertOrUpdateMonitor(monitor)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	o.c.Log.Info("Monitor ", monitor.ID, " added.")
	return nil
}
