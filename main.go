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
	"github.com/kuritsu/spyglass/client"
	"github.com/kuritsu/spyglass/sgc"
	"go.uber.org/fx"
)

func processArgs(cliObj *commands.CommandLineContext) {
	options, err := cli.GetOptions(os.Args[1:])
	cliObj.Log.SetLevel(logr.InfoLevel)
	cliObj.Log.SetFormatter(&LogFormatter{
		ShowDate: err == nil && options.OutputIncludeTimestamps || false,
	})
	if err != nil {
		cliObj.Log.Fatal(err)
		os.Exit(1)
	}
	cliObj.Caller.Init(options.APIAddress)

	cliObj.Log.Println("Setting log level to", options.LogLevel)
	cliObj.Log.SetLevel(options.LogLevelInt)

	runner := options.Action.Apply(cliObj)
	if runner != nil {
		err = runner.Run()
		if err != nil {
			cliObj.Log.Error(err)
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

/*
	All go programs start running from a function called main.
*/
func main() {
	gin.SetMode(gin.ReleaseMode)
	fxlog := logr.New()
	fxlog.SetOutput(ioutil.Discard)

	fmt.Println("spyglass", VERSION)
	app := fx.New(
		fx.Logger(fxlog),
		fx.Provide(
			storage.CreateProviderFromConf,
			func() *logr.Logger {
				return logr.New()
			},
			commands.CreateContext,
			func() sgc.Manager {
				return &sgc.FileManager{}
			},
			client.Create,
		),
		fx.Invoke(processArgs),
	)
	if app.Err() != nil {
		fmt.Println(app.Err().Error())
	}
}
