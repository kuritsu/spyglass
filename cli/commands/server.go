package commands

import (
	"flag"

	"github.com/kuritsu/spyglass/cli/runner"
)

// ServerOptions according to arguments
type ServerOptions struct {
	flagSet *flag.FlagSet
	Address string
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
	fs.StringVar(&result.Address, "a", "127.0.0.1:8010", "API exposed address (Use 0.0.0.0:port for fully open).")
	return &result
}

// Apply the command.
func (o *ServerOptions) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Executing server.")
	s := c.Api.Serve()
	return &runner.Gin{
		Engine:  s,
		Address: o.Address,
	}
}
