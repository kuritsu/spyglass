// Main program
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/cli"
	"github.com/kuritsu/spyglass/cli/commands"
	"github.com/kuritsu/spyglass/sgc"
	"go.uber.org/fx"
)

func processArgs(cliObj *commands.CommandLineContext, logObj *logr.Logger) {
	options, err := cli.GetOptions(os.Args[1:])
	logObj.SetLevel(logr.InfoLevel)
	logObj.SetFormatter(&LogFormatter{
		ShowDate: err == nil && options.OutputIncludeTimestamps || false,
	})
	if err != nil {
		logObj.Fatal(err)
		os.Exit(1)
	}

	logObj.Println("Setting log level to", options.LogLevel)
	logObj.SetLevel(options.LogLevelInt)

	runner := options.Action.Apply(cliObj)
	if runner != nil {
		err = runner.Run()
		if err != nil {
			logObj.Error(err)
			os.Exit(1)
		}
	}
}

// StringListContains tells whether a contains x.
func StringListContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func createLog() *logr.Logger {
	result := logr.New()
	return result
}

/*
	All go programs start running from a function called main.
*/
func main() {
	gin.SetMode(gin.ReleaseMode)
	fxlog := logr.New()
	fxlog.SetOutput(ioutil.Discard)

	fmt.Println("spyglass", VERSION)
	fx.New(
		fx.Logger(fxlog),
		fx.Provide(
			storage.CreateProviderFromConf,
			createLog,
			commands.CreateContext,
			func() sgc.Manager {
				return &sgc.FileManager{}
			},
		),
		fx.Invoke(processArgs),
	)
}
