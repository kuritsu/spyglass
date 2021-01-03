package commands

import (
	"errors"
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/kuritsu/spyglass/cli/commands/mocks"
	"github.com/kuritsu/spyglass/cli/runner"
	"github.com/kuritsu/spyglass/sgc"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
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
	sgcManager := mocks.SgcManagerMock{
		GetFilesResult: []sgc.File{{}, {}, {}},
	}
	r := s.Apply(&CommandLineContext{Db: &testutil.StorageMock{}, Log: mockLog, SgcManager: &sgcManager})

	assert.Nil(t, r)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.DebugLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Processing 3 files")
}

func TestApplyWithNoDirs(t *testing.T) {
	s := ApplyFlags()
	assert.NotNil(t, s)

	f := s.GetFlags()
	assert.NotNil(t, f)

	d := s.Description()
	assert.Contains(t, d, "configuration")

	s.flagSet.Parse([]string{})

	mockLog, _ := test.NewNullLogger()
	r := s.Apply(&CommandLineContext{Db: &testutil.StorageMock{}, Log: mockLog})

	assert.NotNil(t, r)
	exitError := r.(*runner.ExitError)
	assert.Contains(t, exitError.Error.Error(), "Invalid number of directories")
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
	errFunc := s.Apply(&CommandLineContext{Db: &testutil.StorageMock{}, Log: mockLog, SgcManager: &sgcManager})

	assert.NotNil(t, errFunc)
	lastEntry := hook.LastEntry()
	assert.Equal(t, logrus.ErrorLevel, lastEntry.Level)
	assert.Contains(t, lastEntry.Message, "Error reading files")
}
