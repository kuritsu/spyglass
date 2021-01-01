// Main program
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/cli"
	"go.uber.org/fx"
)

func printHelp() {
	fmt.Println("Usage:")
}

func processArgs(apiObj *api.API, cliObj *cli.CommandLine, logObj *logr.Logger) {
	switch os.Args[1] {
	case "server":
		g := apiObj.Serve()
		g.Run()
	default:
		cliObj.Process(os.Args)
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
	result.SetFormatter(&LogFormatter{
		ShowDate: StringListContains(os.Args, "--format-include-timestamps"),
	})
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
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	fx.New(
		fx.Logger(fxlog),
		fx.Provide(
			storage.CreateProviderFromConf,
			createLog,
			api.Create,
			cli.Create,
		),
		fx.Invoke(processArgs),
	)
}
