package sgc

// FileParseError for getting parse errors per file
type FileParseError struct {
	File       *File
	InnerError error
}

// Manager for reading files
type Manager interface {
	GetFiles(string, bool) ([]*File, error)
	ParseConfig() []FileParseError
}

// Error interface implementation
func (p *FileParseError) Error() string {
	return p.InnerError.Error()
}
