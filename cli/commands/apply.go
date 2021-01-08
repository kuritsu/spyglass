package commands

import (
	"errors"
	"flag"
	"fmt"
	"sort"

	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/client"
	"github.com/kuritsu/spyglass/sgc"
)

// ApplyOptions according to arguments
type ApplyOptions struct {
	Recursive         bool
	ForceStatusUpdate bool
	c                 *CommandLineContext
	flagSet           *flag.FlagSet
	fileList          []*sgc.File
}

// GetFlags for the current command.
func (o *ApplyOptions) GetFlags() *flag.FlagSet {
	return o.flagSet
}

// Description for the current command.
func (o *ApplyOptions) Description() string {
	return "Apply the configuration contained in *.sgc files in the specified directory."
}

// ApplyFlags obtains flags for apply action
func ApplyFlags() *ApplyOptions {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	result := ApplyOptions{flagSet: fs}
	fs.BoolVar(&result.Recursive, "r", false, "Scan specified paths recursively for config files.")
	fs.BoolVar(&result.ForceStatusUpdate, "fsu", false, "Forces the status update for targets.")
	fs.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  spyglass apply [flags] path1 [path2] ...")
		fmt.Println("\nFlags:")
		fs.PrintDefaults()
	}
	return &result
}

// Apply the configuration in the given directory.
func (o *ApplyOptions) Apply(c *CommandLineContext) runner.Runner {
	o.c = c
	c.Log.Debug("Executing apply.")
	var err runner.Runner
	o.fileList, err = o.getFiles()
	if err != nil {
		return err
	}
	c.Log.Debug("Processing ", len(o.fileList), " files...")
	errors := c.SgcManager.ParseConfig()
	if errors != nil {
		for _, f := range errors {
			c.Log.Errorln(f.Error())
		}
		return &runner.ExitError{}
	}
	err = o.applyConfigs()
	if err != nil {
		return err
	}
	return nil
}

func (o *ApplyOptions) getFiles() ([]*sgc.File, runner.Runner) {
	paths := o.flagSet.Args()
	if len(paths) == 0 {
		return nil, &runner.ExitError{
			Error:   errors.New("A path is required"),
			FlagSet: o.flagSet,
			Logger:  o.c.Log,
		}
	}
	fileList := []*sgc.File{}
	for _, d := range paths {
		files, err := o.c.SgcManager.GetFiles(d, o.Recursive)
		if err != nil {
			o.c.Log.Error(err)
			return nil, &runner.ExitError{
				Error:  err,
				Logger: o.c.Log,
			}
		}
		fileList = append(fileList, files...)
	}
	return fileList, nil
}

func (o *ApplyOptions) applyConfigs() runner.Runner {
	allMonitors := []*types.Monitor{}
	allTargets := []*types.Target{}
	for _, f := range o.fileList {
		if f.Config.Monitors != nil {
			allMonitors = append(allMonitors, f.Config.Monitors...)
		}
		if f.Config.Targets != nil {
			allTargets = append(allTargets, f.Config.Targets...)
		}
	}
	if err := o.applyAllMonitors(o.c.Caller, allMonitors); err != nil {
		o.c.Log.Debug("Error applying monitors config.")
		return &runner.ExitError{Error: err}
	}
	if err := o.applyAllTargets(o.c.Caller, types.TargetList(allTargets)); err != nil {
		o.c.Log.Debug("Error applying targets config.")
		return &runner.ExitError{Error: err}
	}
	return nil
}

func (o *ApplyOptions) applyAllMonitors(api client.APICaller, monitors []*types.Monitor) error {
	for _, m := range monitors {
		if err := api.InsertOrUpdateMonitor(m); err != nil {
			return err
		}
	}
	return nil
}

func (o *ApplyOptions) applyAllTargets(api client.APICaller, targets types.TargetList) error {
	sort.Sort(targets)
	for _, m := range targets {
		if err := api.InsertOrUpdateTarget(m, o.ForceStatusUpdate); err != nil {
			return err
		}
	}
	return nil
}
