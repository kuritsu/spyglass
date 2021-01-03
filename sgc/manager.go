package sgc

// Manager for reading files
type Manager interface {
	GetFiles(string, bool) ([]File, error)
}
