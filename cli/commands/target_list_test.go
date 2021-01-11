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

func TestTargetListAction_GetFlags(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetListActionFlags(parentFs)
	fls := action.GetFlags()
	assert.NotNil(t, fls)
	fls.Usage()
	assert.NotNil(t, action.Description())
}

func TestTargetListActionListTargetsError(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetListActionFlags(parentFs)
	assert.NotNil(t, action)

	mockLog, _ := test.NewNullLogger()
	mockLog.SetLevel(logrus.DebugLevel)
	caller := client.CallerMock{}
	caller.On("ListTargets", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Connection error"))
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
