package cli

import (
	"fmt"
	"os"

	"github.com/kuritsu/spyglass/api/storage"
)

// CommandLine processing with API calls
type CommandLine struct {
	db storage.Provider
}

// Create an instance of the CLI object
func Create(db storage.Provider) *CommandLine {
	result := CommandLine{db}
	return &result
}

// Process the command line arguments
func (c *CommandLine) Process() {
	switch os.Args[1] {
	case "apply":
		fmt.Println("apply")
	}
}
