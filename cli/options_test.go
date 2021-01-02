package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts, err := GetOptions([]string{})
	assert.NotNil(t, opts)
	assert.Nil(t, err)
}
