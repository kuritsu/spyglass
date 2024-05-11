package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/kuritsu/spyglass/cli/runner"
)

// UserTokenAction represents the list action of targets.
type UserTokenAction struct {
	flagSet *flag.FlagSet
	c       *CommandLineContext
	email   string
	hours   int
}

// GetFlags for the current command.
func (o *UserTokenAction) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *UserTokenAction) Description() string {
	return "Create an API token for the given user."
}

// UserTokenActionFlags obtains flags for user token action
func UserTokenActionFlags(parentFs *flag.FlagSet) *UserTokenAction {
	result := UserTokenAction{}
	result.flagSet = flag.NewFlagSet("token", flag.ContinueOnError)
	result.flagSet.StringVar(&result.email, "u", "", "User email.")
	result.flagSet.IntVar(&result.hours, "h", 24, "Hours to expire. Default: 24. Cannot be bigger than 8760 (1 year) nor less than 1.")
	result.flagSet.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass user [global-flags] token -u user@email.com -h 48")
		fmt.Println("\nFlags:")
		result.flagSet.PrintDefaults()
		fmt.Println("\nGlobal Flags:")
		parentFs.PrintDefaults()
	}
	return &result
}

// Apply the current action.
func (o *UserTokenAction) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Apply user token...")
	if o.hours < 1 || o.hours > 8760 {
		return &runner.ExitError{Error: fmt.Errorf("Invalid hours argument. Must be between 1 and 8760."), Logger: c.Log}
	}
	expiration := time.Now().Add(time.Hour * time.Duration(o.hours))
	token, err := c.Caller.CreateUserToken(o.email, expiration)
	if err != nil {
		return &runner.ExitError{Error: err, Logger: c.Log}
	}
	c.Log.Info(token)
	return nil
}
