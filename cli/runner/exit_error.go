package runner

import (
	"flag"
	"os"

	logr "github.com/sirupsen/logrus"
)

// ExitError runner.
type ExitError struct {
	Error   error
	FlagSet *flag.FlagSet
	Logger  *logr.Logger
}

// Run engine.
func (e *ExitError) Run() error {
	if e.FlagSet == nil {
		return e.Error
	}
	e.Logger.Error(e.Error)
	e.FlagSet.Usage()
	os.Exit(1)
	return nil
}
