// Main program
package main

import (
	"os"

	"github.com/kuritsu/spyglass/api"
)

/*
	All go programs start running from a function called main.
*/
func main() {
	switch os.Args[1] {
	case "server":
		s := api.Serve()
		s.Run()
	}
}
