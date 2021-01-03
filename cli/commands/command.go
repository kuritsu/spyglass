package commands

import (
	"flag"

	"github.com/kuritsu/spyglass/cli/runner"
)

// Command represents a CLI command
type Command interface {
	GetFlags() *flag.FlagSet
	Apply(*CommandLineContext) runner.Runner
	Description() string
}
