package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTargetParentByIDNoParent(t *testing.T) {
	s := "target"
	res := GetTargetParentByID(s)
	assert.Equal(t, "", res)
}

func TestGetTargetParentByIDParent(t *testing.T) {
	s := "parent.target"
	res := GetTargetParentByID(s)
	assert.Equal(t, "parent", res)
}

func TestGetTargetParentByIDBigParent(t *testing.T) {
	s := "parent1.parent-2.parent_3.hello"
	res := GetTargetParentByID(s)
	assert.Equal(t, "parent1.parent-2.parent_3", res)
}
