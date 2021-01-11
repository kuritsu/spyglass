package commands

import (
	"errors"
	"flag"
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/client"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTargetUpdateStatusAction_GetFlags(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	fls := action.GetFlags()
	assert.NotNil(t, fls)
	fls.Usage()
	assert.NotNil(t, action.Description())
}

func TestTargetUpdateStatusActionWithConnectionError(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	assert.NotNil(t, action)
	action.flagSet.Parse([]string{"-id", "target", "-s", "50"})

	mockLog, _ := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	caller.On("UpdateTargetStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("Connection error"))
	result := action.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &caller,
		SgcManager: nil,
	})

	assert.NotNil(t, result)
	assert.IsType(t, &runner.ExitError{}, result)
	exitError := result.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "Connection error")
}

func TestTargetUpdateStatusActionSuccess(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	action.flagSet.Parse([]string{"-id", "target", "-s", "50"})
	assert.NotNil(t, action)

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	caller.On("UpdateTargetStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	result := action.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &caller,
		SgcManager: nil,
	})

	assert.Nil(t, result)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Status updated")
}

func TestTargetUpdateStatusActionNoIdError(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	assert.NotNil(t, action)
	action.flagSet.Parse([]string{})

	mockLog, _ := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	result := action.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &caller,
		SgcManager: nil,
	})

	assert.NotNil(t, result)
	assert.IsType(t, &runner.ExitError{}, result)
	exitError := result.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "id (target ID) flag is required")
}

func TestTargetUpdateStatusActionNoStatusError(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	assert.NotNil(t, action)
	action.flagSet.Parse([]string{"-id", "target"})

	mockLog, _ := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	result := action.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &caller,
		SgcManager: nil,
	})

	assert.NotNil(t, result)
	assert.IsType(t, &runner.ExitError{}, result)
	exitError := result.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "s (status) flag is required")
}
