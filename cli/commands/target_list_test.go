package commands

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetListAction_GetFlags(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetListActionFlags(parentFs)
	fls := action.GetFlags()
	assert.NotNil(t, fls)
	fls.Usage()
	assert.NotNil(t, action.Description())
}
