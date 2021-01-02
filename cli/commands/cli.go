package commands

import (
	"github.com/kuritsu/spyglass/api/storage"
	logr "github.com/sirupsen/logrus"
)

// CommandLineContext processing with API calls
type CommandLineContext struct {
	Db  storage.Provider
	Log *logr.Logger
}

// CreateContext an instance of the CLI object
func CreateContext(db storage.Provider, log *logr.Logger) *CommandLineContext {
	result := CommandLineContext{db, log}
	result.Log.Debug("Created CommandLine instance.")
	return &result
}
