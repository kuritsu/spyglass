package api

import (
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApiServe(t *testing.T) {
	s := Create(&testutil.StorageMock{}, logrus.New())
	s.Serve()
	assert.NotNil(t, s)
}

func TestApiServeWithDebug(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s := Create(&testutil.StorageMock{}, l)
	s.Serve()
	assert.NotNil(t, s)
}
