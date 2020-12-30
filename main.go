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

func processArgs(apiObj *api.API, cliObj *cli.CommandLine) {
	switch os.Args[1] {
	case "server":
		apiObj.Serve()
	default:
		cliObj.Process()
	}
}

/*
	All go programs start running from a function called main.
*/
func main() {
	gin.SetMode(gin.ReleaseMode)
	std := logr.StandardLogger()
	std.SetOutput(ioutil.Discard)

	fmt.Println("spyglass", VERSION)
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	fx.New(
		fx.Logger(std),
		fx.Provide(
			storage.CreateProviderFromConf,
			api.Create,
			cli.Create,
		),
		fx.Invoke(processArgs),
	)
}
