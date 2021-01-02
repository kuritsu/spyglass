package commands

import (
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestCommandLine(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)
	c := CreateContext(&testutil.StorageMock{}, logger)

	assert.NotNil(t, c)
	assert.NotNil(t, hook.AllEntries())
	assert.Len(t, hook.AllEntries(), 1)
	assert.Contains(t, hook.AllEntries()[0].Message, "Created CommandLine instance.")
}

func TestCommandLineApply(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)
	c := CreateContext(&testutil.StorageMock{}, logger)

	c.Process(&Options{Action: "apply"})

	entries := hook.AllEntries()
	assert.Contains(t, entries[len(entries)-1].Message, "Executing apply")
}
