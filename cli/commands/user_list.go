package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kuritsu/spyglass/cli/runner"
)

// UserListAction represents the list action of targets.
type UserListAction struct {
	flagSet             *flag.FlagSet
	c                   *CommandLineContext
	format              string
	pageIndex, pageSize int
}

// GetFlags for the current command.
func (o *UserListAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *UserListAction) Description() string {
	return "List users."
}

// UserListActionFlags obtains flags for target list action
func UserListActionFlags(parentFs *flag.FlagSet) *UserListAction {
	result := UserListAction{}
	result.flagSet = flag.NewFlagSet("list", flag.ContinueOnError)
	result.flagSet.StringVar(&result.format, "o", "json", "Output format. Can be json.")
	result.flagSet.IntVar(&result.pageIndex, "pi", 0, "Page index.")
	result.flagSet.IntVar(&result.pageSize, "ps", 100, "Page size.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass user [global-flags] list -pi 1 -ps 50")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *UserListAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply user list...")
	users, err := c.Caller.ListUsers(o.pageIndex, o.pageSize)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	switch o.format {
	case "json":
		c.Log.Debug("Displaying users as json...")
		jsonBytes, _ := json.Marshal(users)
		var out bytes.Buffer
		json.Indent(&out, jsonBytes, "", "  ")
		out.WriteTo(os.Stdout)
	}
	return nil
}
