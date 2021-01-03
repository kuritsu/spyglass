package sgc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	m := FileManager{}
	files, err := m.GetFiles("./e2e", true)
	assert.NotNil(t, files)
	assert.Nil(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, "file.sgc", files[0].File.Name())
	for _, f := range files {
		assert.NotEqual(t, "file.sgcf", f.File.Name())
	}
}

func TestGetFilesWithInvalidDir(t *testing.T) {
	m := FileManager{}
	files, err := m.GetFiles("./e2exx", true)
	assert.Nil(t, files)
	assert.NotNil(t, err)
}
