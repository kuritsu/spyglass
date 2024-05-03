package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuritsu/spyglass/cli/runner"
)

// ServerOptions according to arguments
type LoginOptions struct {
	flagSet  *flag.FlagSet
	user     string
	password string
}

// GetFlags for the current command.
func (o *LoginOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *LoginOptions) Description() string {
	return "Log in the API server."
}

// ServerFlags obtains flags for apply action
func LoginFlags() *LoginOptions {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	result := LoginOptions{flagSet: fs}
	fs.StringVar(&result.user, "u", "", "User email.")
	fs.StringVar(&result.password, "p", "", "User password.")
	return &result
}

// Apply the command.
func (o *LoginOptions) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Logging in...")
	token, err := c.Caller.Login(o.user, o.password)
	if err != nil {
		c.Log.Error(err)
		return nil
	}
	homedir, _ := os.UserHomeDir()
	fname := filepath.Join(homedir, ".spyglass.token")
	err = os.WriteFile(fname, []byte(fmt.Sprintf("%s:%s", o.user, token)), 0660)
	if err != nil {
		c.Log.Error(err)
	}
	return nil
}
