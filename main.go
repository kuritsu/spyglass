// Main program
package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api"
	"github.com/kuritsu/spyglass/api/storage"
	"go.uber.org/fx"
)

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
	fx.New(
		fx.Provide(
			storage.CreateProviderFromConf,
			api.Serve,
		),
		fx.Invoke(processArgs),
	)
}
