package commands

import (
	"embed"

	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/client"
	logr "github.com/sirupsen/logrus"
)

// CommandLineContext processing with API calls
type CommandLineContext struct {
	Db     storage.Provider
	Log    *logr.Logger
	Caller client.APICaller
	res    embed.FS
}

// CreateContext an instance of the CLI object
func CreateContext(db storage.Provider, log *logr.Logger,
	api client.APICaller, res embed.FS) *CommandLineContext {
	result := CommandLineContext{db, log, api, res}
	result.Log.Debug("Created CommandLine instance.")
	return &result
}
