package commands

import "flag"

// Command represents a CLI command
type Command interface {
	GetFlags() *flag.FlagSet
	Apply(*CommandLineContext) func(...string) error
	Description() string
}
