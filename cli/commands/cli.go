package commands

import (
	"embed"

	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/client"
	"github.com/kuritsu/spyglass/scheduler"
	logr "github.com/sirupsen/logrus"
)

// CommandLineContext processing with API calls
type CommandLineContext struct {
	Db     storage.Provider
	Log    *logr.Logger
	Caller client.APICaller
	res    embed.FS
	Api    *api.API
	Sch    scheduler.Scheduler
}

// CreateContext an instance of the CLI object
func CreateContext(db storage.Provider, log *logr.Logger,
	apiClient client.APICaller, res embed.FS, apiServer *api.API, sch scheduler.Scheduler) *CommandLineContext {
	result := CommandLineContext{db, log, apiClient, res, apiServer, sch}
	result.Log.Debug("Created CommandLine instance.")
	return &result
}
