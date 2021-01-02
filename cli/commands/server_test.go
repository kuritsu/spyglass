package commands

import (
	"testing"

	"github.com/kuritsu/spyglass/api/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	s := ServerFlags()
	assert.NotNil(t, s)

	f := s.GetFlags()
	assert.NotNil(t, f)

	d := s.Description()
	assert.Contains(t, d, "API")

	s.Apply(&CommandLineContext{Db: &testutil.StorageMock{}, Log: logrus.New()})
}
