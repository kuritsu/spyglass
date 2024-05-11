package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kuritsu/spyglass/cli/runner"
)

// RoleRevokeAction represents the revoke action of roles.
type RoleRevokeAction struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	role    string
	users   string
}

// GetFlags for the current command.
func (o *RoleRevokeAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *RoleRevokeAction) Description() string {
	return "Revokes a role assignment to a list of users."
}

// RoleRevokeActionFlags obtains flags for target list action
func RoleRevokeActionFlags(parentFs *flag.FlagSet) *RoleRevokeAction {
	result := RoleRevokeAction{}
	result.flagSet = flag.NewFlagSet("revoke", flag.ContinueOnError)
	result.flagSet.StringVar(&result.role, "r", "", "Role name.")
	result.flagSet.StringVar(&result.users, "u", "", "Comma separated user emails.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass role [global-flags] revoke -r roleName -u user1@email.com,user2@email.com")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *RoleRevokeAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply role revoke...")
	usersRemove := strings.Split(o.users, ",")
	err := c.Caller.UpdateRole(o.role, nil, usersRemove)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	return nil
}
