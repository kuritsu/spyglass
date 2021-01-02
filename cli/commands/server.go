package commands

import (
	"flag"

	"github.com/kuritsu/spyglass/api"
)

// ServerOptions according to arguments
type ServerOptions struct {
	flagSet *flag.FlagSet
}

// Apply the command.
func (o *ServerOptions) Apply(c *CommandLineContext) {
	c.Log.Debug("Executing server.")
	apiObj := api.Create(c.Db, c.Log)
	s := apiObj.Serve()
	s.Run()
}

// GetFlags for the current command.
func (o *ServerOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *ServerOptions) Description() string {
	return "Execute the API server."
}

// ServerFlags obtains flags for apply action
func ServerFlags() *ServerOptions {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	result := ServerOptions{flagSet: fs}
	return &result
}
