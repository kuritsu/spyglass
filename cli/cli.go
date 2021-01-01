package cli

import (
	"github.com/kuritsu/spyglass/api/storage"
	logr "github.com/sirupsen/logrus"
)

// CommandLine processing with API calls
type CommandLine struct {
	db  storage.Provider
	log *logr.Logger
}

// Create an instance of the CLI object
func Create(db storage.Provider, log *logr.Logger) *CommandLine {
	result := CommandLine{db, log}
	result.log.Debug("Created CommandLine instance.")
	return &result
}

// Process the command line arguments
func (c *CommandLine) Process(args []string) {
	switch args[1] {
	case "apply":
		c.Apply(args[2], ApplyOptions{})
	}
}
