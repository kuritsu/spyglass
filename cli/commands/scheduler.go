package commands

import (
	"flag"

	"github.com/kuritsu/spyglass/cli/runner"
)

// SchedulerOptions according to arguments
type SchedulerOptions struct {
	flagSet *flag.FlagSet
	Label   string
}

// GetFlags for the current command.
func (o *SchedulerOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *SchedulerOptions) Description() string {
	return "Execute the Scheduler server."
}

// SchedulerOptionsFlags obtains flags for apply action
func SchedulerOptionsFlags() *SchedulerOptions {
	fs := flag.NewFlagSet("scheduler", flag.ContinueOnError)
	result := SchedulerOptions{flagSet: fs}
	fs.StringVar(&result.Label, "l", "", "Scheduler label (to distiguish from a group of scheduler instances).")
	return &result
}

// Apply the command.
func (o *SchedulerOptions) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Executing Scheduler process.")
	return nil
}
