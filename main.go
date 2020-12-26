// Main program
package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/storage"
	"go.uber.org/fx"
)

func printHelp() {
	fmt.Println("Usage:")
}

func processArgs(s *gin.Engine) {
	switch os.Args[1] {
	case "server":
		s.Run()
	}
}

/*
	All go programs start running from a function called main.
*/
func main() {
	fmt.Println("spyglass", VERSION)
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	fx.New(
		fx.Provide(
			storage.CreateProviderFromConf,
			api.Serve,
		),
		fx.Invoke(processArgs),
	)
}
