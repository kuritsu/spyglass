package api

// ErrorWithCode for http results
type ErrorWithCode struct {
	Message string
	Code    int
}

// Error message
func (e *ErrorWithCode) Error() string {
	return e.Message
}
