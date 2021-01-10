package commands

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetUpdateStatusAction_GetFlags(t *testing.T) {
	parentFs := flag.NewFlagSet("target", flag.ContinueOnError)
	action := TargetUpdateStatusActionFlags(parentFs)
	fls := action.GetFlags()
	assert.NotNil(t, fls)
	fls.Usage()
	assert.NotNil(t, action.Description())
}
