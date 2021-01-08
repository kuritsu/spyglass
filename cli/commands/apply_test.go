package commands

import (
	"errors"
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/kuritsu/spyglass/api/types"
	"github.com/kuritsu/spyglass/cli/commands/mocks"
	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/client"
	"github.com/kuritsu/spyglass/sgc"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApply(t *testing.T) {
	s := ApplyFlags()
	assert.NotNil(t, s)

	f := s.GetFlags()
	assert.NotNil(t, f)

	d := s.Description()
	assert.Contains(t, d, "configuration")

	s.flagSet.Parse([]string{"."})

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	sgcManager := getGoodSgcManager()
	caller := &client.CallerMock{}
	caller.On("InsertOrUpdateMonitor", mock.Anything).Return(nil)
	caller.On("InsertOrUpdateTarget", mock.Anything, false).Return(nil)
	r := s.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		SgcManager: sgcManager,
		Caller:     caller,
	})

	assert.Nil(t, r)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Processing 3 files")
}

func TestApplyWithNoDirs(t *testing.T) {
	s := ApplyFlags()
	s.GetFlags()
	d := s.Description()
	assert.Contains(t, d, "configuration")

	s.flagSet.Parse([]string{})

	mockLog, _ := test.NewNullLogger()
	r := s.Apply(&CommandLineContext{
		Db:     &testutil.StorageMock{},
		Log:    mockLog,
		Caller: &client.CallerMock{},
	})

	assert.NotNil(t, r)
	exitError := r.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "A path is required")
	assert.NotNil(t, exitError.FlagSet)
	exitError.FlagSet.Usage()
}

func TestApplyGetFilesError(t *testing.T) {
	s := ApplyFlags()
	s.flagSet.Parse([]string{"."})

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	sgcManager := mocks.SgcManagerMock{
		GetFilesError: errors.New("Error reading files"),
	}
	errFunc := s.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &client.CallerMock{},
		SgcManager: &sgcManager,
	})

	assert.NotNil(t, errFunc)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.ErrorLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Error reading files")
}

func TestApplyParseFilesError(t *testing.T) {
	s := ApplyFlags()
	s.flagSet.Parse([]string{"."})

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	sgcManager := mocks.SgcManagerMock{
		ParseConfigResult: []sgc.FileParseError{{File: nil, InnerError: errors.New("Invalid file")}},
	}
	errFunc := s.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     &client.CallerMock{},
		SgcManager: &sgcManager,
	})

	assert.NotNil(t, errFunc)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.ErrorLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Invalid file")
}

func TestApplyApplyMonitorConfigsError(t *testing.T) {
	s := ApplyFlags()
	s.flagSet.Parse([]string{"."})

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	sgcManager := getGoodSgcManager()
	caller := &client.CallerMock{}
	caller.On("InsertOrUpdateMonitor", mock.Anything).Return(errors.New("Invalid monitor"))
	errFunc := s.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     caller,
		SgcManager: sgcManager,
	})

	assert.NotNil(t, errFunc)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Error applying monitors config")
}

func TestApplyApplyTargetConfigsError(t *testing.T) {
	s := ApplyFlags()
	s.flagSet.Parse([]string{"."})

	mockLog, hook := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	sgcManager := getGoodSgcManager()
	caller := &client.CallerMock{}
	caller.On("InsertOrUpdateMonitor", mock.Anything).Return(nil)
	caller.On("InsertOrUpdateTarget", mock.Anything, false).Return(errors.New("Invalid target"))
	errFunc := s.Apply(&CommandLineContext{
		Db:         &testutil.StorageMock{},
		Log:        mockLog,
		Caller:     caller,
		SgcManager: sgcManager,
	})

	assert.NotNil(t, errFunc)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Error applying targets config")
}

func getGoodSgcManager() *mocks.SgcManagerMock {
	return &mocks.SgcManagerMock{
		GetFilesResult: []*sgc.File{
			{Config: &sgc.FileConfig{Monitors: []*types.Monitor{{ID: "m1"}}}},
			{Config: &sgc.FileConfig{Targets: []*types.Target{{ID: "m1"}}}},
			{Config: &sgc.FileConfig{Monitors: []*types.Monitor{{ID: "m1"}}, Targets: []*types.Target{{ID: "m1"}}}}},
	}
}
