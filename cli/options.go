package cli

import (
	"flag"
	"fmt"
	"os"

	logr "github.com/sirupsen/logrus"

	"github.com/kuritsu/spyglass/cli/commands"
)

// Options for the command line
type Options struct {
	Action                  commands.Command
	LogLevel                string
	LogLevelInt             logr.Level
	OutputIncludeTimestamps bool
	Help                    bool
	APIAddress              string
}

var logLevelNames = map[string]logr.Level{
	"FATAL": logr.FatalLevel,
	"ERROR": logr.ErrorLevel,
	"WARN":  logr.WarnLevel,
	"INFO":  logr.InfoLevel,
	"DEBUG": logr.DebugLevel,
}

// CommandList list of commands.
var CommandList = map[string]commands.Command{
	"login":  commands.LoginFlags(),
	"server": commands.ServerFlags(),
	"target": commands.TargetFlags(),
	"ui":     commands.UIOptionsFlags(),
}

func defineGlobalFlags(fs *flag.FlagSet, opts *Options) {
	fs.StringVar(&opts.APIAddress,
		"api", "http://localhost:8010", "Spyglass API address. Environment: SPYGLASS_API.")
	fs.BoolVar(&opts.OutputIncludeTimestamps,
		"ts", false, "Include timestamps on output.")
	fs.StringVar(&opts.LogLevel,
		"v", "INFO", "Verbosity. Possible values: FATAL, ERROR, WARN, INFO, DEBUG.")
	fs.BoolVar(&opts.Help, "h", false, "Help.")
}

func printGlobalHelp(fs *flag.FlagSet) {
	fmt.Println("Usage:")
	fmt.Println("  spyglass <command> [global-flags]")
	fmt.Println("\nCommands:")
	for k, v := range CommandList {
		fmt.Println("  ", k)
		fmt.Println("    ", v.Description())
	}
	fmt.Println("\nGlobal flags:")
	fs.PrintDefaults()
}

// GetOptions parses command line to get the options
func GetOptions(args []string) (*Options, error) {
	result := Options{}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	defineGlobalFlags(fs, &result)
	err := fs.Parse(args)
	if err != nil && err != flag.ErrHelp {
		return nil, err
	}
	if len(args) == 0 {
		printGlobalHelp(fs)
		os.Exit(0)
	}

	var ok bool
	result.Action, ok = CommandList[args[0]]
	if !ok {
		printGlobalHelp(fs)
		os.Exit(0)
	}

	actionFlags := result.Action.GetFlags()
	defineGlobalFlags(actionFlags, &result)
	actionFlags.Parse(args[1:])
	if result.Help {
		result.Action.GetFlags().Usage()
		os.Exit(0)
	}

	level, ok := logLevelNames[result.LogLevel]
	if !ok {
		level = logr.InfoLevel
	}
	result.LogLevelInt = level

	return &result, nil
}
