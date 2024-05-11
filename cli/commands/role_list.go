package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kuritsu/spyglass/cli/runner"
)

// RoleListAction represents the list action of targets.
type RoleListAction struct {
	flagSet             *flag.FlagSet
	c                   *CommandLineContext
	format              string
	pageIndex, pageSize int
}

// GetFlags for the current command.
func (o *RoleListAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *RoleListAction) Description() string {
	return "List roles."
}

// RoleListActionFlags obtains flags for target list action
func RoleListActionFlags(parentFs *flag.FlagSet) *RoleListAction {
	result := RoleListAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.format, "o", "json", "Output format. Can be json.")
	result.flagSet.IntVar(&result.pageIndex, "pi", 0, "Page index.")
	result.flagSet.IntVar(&result.pageSize, "ps", 100, "Page size.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass role [global-flags] list -pi 1 -ps 50")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *RoleListAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply role list...")
	roles, err := c.Caller.ListRoles(o.pageIndex, o.pageSize)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	switch o.format {
	case "json":
		c.Log.Debug("Displaying roles as json...")
		jsonBytes, _ := json.Marshal(roles)
		var out bytes.Buffer
		json.Indent(&out, jsonBytes, "", "  ")
		out.WriteTo(os.Stdout)
	}
	return nil
}
