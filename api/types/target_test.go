package types

import (
	"sort"
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

func TestGetIDForRegex(t *testing.T) {
	s := "parent1.parent-2.Parent_3.hello"
	res := GetIDForRegex(s)
	assert.Equal(t, `parent1\.parent\-2\.parent_3\.hello`, res)
}

func TestTargetListSorting(t *testing.T) {
	targetList := TargetList{{ID: "z"}, {ID: "X"}, {ID: "Y"}}
	sort.Sort(targetList)
	assert.Equal(t, targetList[0].ID, "X")
	assert.Equal(t, targetList[1].ID, "Y")
	assert.Equal(t, targetList[2].ID, "z")
}
