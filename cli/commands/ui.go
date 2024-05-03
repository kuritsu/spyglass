package commands

import (
	"flag"
	"io/fs"
	"net/http"

	"github.com/kuritsu/spyglass/cli/runner"
)

// UIOptions according to arguments
type UIOptions struct {
	flagSet *flag.FlagSet
	Address string
}

// GetFlags for the current command.
func (o *UIOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *UIOptions) Description() string {
	return "Execute the UI server."
}

// UIOptionsFlags obtains flags for apply action
func UIOptionsFlags() *UIOptions {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	result := UIOptions{flagSet: fs}
	fs.StringVar(&result.Address, "a", "127.0.0.1:8080", "API exposed address (Use 0.0.0.0:port for fully open).")
	return &result
}

// Apply the command.
func (o *UIOptions) Apply(c *CommandLineContext) runner.Runner {
	c.Log.Debug("Executing UI server.")
	c.Log.Info("Listening on ", o.Address, "...")
	content, _ := fs.Sub(c.res, "ui/dist")
	fserver := http.FileServer(http.FS(content))
	http.Handle("/target", http.StripPrefix("/target", fserver))
	http.Handle("/", fserver)
	err := http.ListenAndServe(o.Address, nil)
	if err != nil {
		panic(err)
	}
	return nil
}
