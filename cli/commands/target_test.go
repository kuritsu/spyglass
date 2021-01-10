package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
