package commands

import (
	"flag"
	"slices"

	"github.com/kuritsu/spyglass/cli/runner"
)

// Command represents a CLI command
type Command interface {
	GetFlags() *flag.FlagSet
	Apply(*CommandLineContext) runner.Runner
	Description() string
}

func GetSortedKeyList(cmdMap map[string]Command) []string {
	commands := make([]string, 0, len(cmdMap))
	for k := range cmdMap {
		commands = append(commands, k)
	}
	slices.Sort(commands)
	return commands
}
