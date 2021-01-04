package mocks

import "github.com/kuritsu/spyglass/sgc"

// SgcManagerMock mock
type SgcManagerMock struct {
	GetFilesResult    []*sgc.File
	GetFilesError     error
	ParseConfigResult []sgc.FileParseError
}

// GetFiles mock
func (m *SgcManagerMock) GetFiles(dir string, recursive bool) ([]*sgc.File, error) {
	return m.GetFilesResult, m.GetFilesError
}

// ParseConfig from file list
func (m *SgcManagerMock) ParseConfig() []sgc.FileParseError {
	return m.ParseConfigResult
}
