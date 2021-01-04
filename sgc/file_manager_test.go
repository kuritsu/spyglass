package sgc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	m := FileManager{}
	files, err := m.GetFiles("./e2e", true)
	assert.NotNil(t, files)
	assert.Nil(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, "file.hcl", files[0].File.Name())
	for _, f := range files {
		assert.NotEqual(t, "file.hclf", f.File.Name())
	}
}

func TestGetFilesWithInvalidDir(t *testing.T) {
	m := FileManager{}
	files, err := m.GetFiles("./e2exx", true)
	assert.Nil(t, files)
	assert.NotNil(t, err)
}

func TestParseConfigNoFilesRead(t *testing.T) {
	m := FileManager{}
	errors := m.ParseConfig()
	assert.NotNil(t, errors)
	assert.Len(t, errors, 1)
	assert.Equal(t, "No files read", errors[0].Error())
}

func TestParseConfigWithFiles(t *testing.T) {
	m := FileManager{}
	files, err := m.GetFiles("./e2e", false)
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	errors := m.ParseConfig()
	fmt.Print(errors)
	assert.Nil(t, errors)
	assert.NotNil(t, files[0].Config)
}
