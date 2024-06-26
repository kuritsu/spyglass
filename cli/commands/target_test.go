package commands

import (
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/client"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTargetFlags(t *testing.T) {
	got := TargetFlags()
	fs := got.GetFlags()
	assert.NotNil(t, fs)
	assert.NotEmpty(t, got.Description())
	fs.Usage()
	fs.Parse([]string{"-h", "list"})
	fs.Usage()
}

func TestTargetFlagsWithNoAction(t *testing.T) {
	got := TargetFlags()

	got.flagSet.Parse([]string{})
	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := &client.CallerMock{}

	result := got.Apply(&CommandLineContext{
		Db:     &testutil.StorageMock{},
		Log:    mockLog,
		Caller: caller,
	})

	assert.NotNil(t, result)
	assert.IsType(t, &runner.ExitError{}, result)
	exitError := result.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "An action is required")
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Doing target")
}

func TestTargetFlagsWithUnsupportedAction(t *testing.T) {
	got := TargetFlags()

	got.flagSet.Parse([]string{"unsupported"})
	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	result := got.Apply(&CommandLineContext{
		Db:     &testutil.StorageMock{},
		Log:    mockLog,
		Caller: &client.CallerMock{},
	})

	assert.NotNil(t, result)
	assert.IsType(t, &runner.ExitError{}, result)
	exitError := result.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "Action not supported")
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Doing target")
}

func TestTargetFlagsWithListAction(t *testing.T) {
	got := TargetFlags()

	got.flagSet.Parse([]string{"list"})
	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	caller.On("ListTargets", mock.Anything, mock.Anything, mock.Anything).Return([]*types.Target{
		{ID: "target"},
	}, nil)
	result := got.Apply(&CommandLineContext{
		Db:     &testutil.StorageMock{},
		Log:    mockLog,
		Caller: &caller,
	})

	assert.Nil(t, result)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Displaying targets as json")
}
