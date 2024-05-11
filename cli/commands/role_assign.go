package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kuritsu/spyglass/cli/runner"
)

// RoleAssignAction represents the list action of targets.
type RoleAssignAction struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	role    string
	users   string
}

// GetFlags for the current command.
func (o *RoleAssignAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *RoleAssignAction) Description() string {
	return "Assigns a role to a list of users."
}

// RoleAssignActionFlags obtains flags for target list action
func RoleAssignActionFlags(parentFs *flag.FlagSet) *RoleAssignAction {
	result := RoleAssignAction{}
	result.flagSet = flag.NewFlagSet("assign", flag.ContinueOnError)
	result.flagSet.StringVar(&result.role, "r", "", "Role name.")
	result.flagSet.StringVar(&result.users, "u", "", "Comma separated user emails.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass role [global-flags] assign -r roleName -u user1@email.com,user2@email.com")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *RoleAssignAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply role assign...")
	usersAdd := strings.Split(o.users, ",")
	err := c.Caller.UpdateRole(o.role, usersAdd, nil)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	return nil
}
