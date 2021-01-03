package mocks

import "github.com/kuritsu/spyglass/sgc"

// SgcManagerMock mock
type SgcManagerMock struct {
	GetFilesResult []sgc.File
	GetFilesError  error
}

// GetFiles mock
func (m *SgcManagerMock) GetFiles(dir string, recursive bool) ([]sgc.File, error) {
	return m.GetFilesResult, m.GetFilesError
}
