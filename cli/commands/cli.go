package commands

import (
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/client"
	"github.com/kuritsu/spyglass/sgc"
	logr "github.com/sirupsen/logrus"
)

// CommandLineContext processing with API calls
type CommandLineContext struct {
	Db         storage.Provider
	Log        *logr.Logger
	SgcManager sgc.Manager
	Caller     client.APICaller
}

// CreateContext an instance of the CLI object
func CreateContext(db storage.Provider, log *logr.Logger, fmgr sgc.Manager,
	api client.APICaller) *CommandLineContext {
	result := CommandLineContext{db, log, fmgr, api}
	result.Log.Debug("Created CommandLine instance.")
	return &result
}
