package cli

// ExitError carries an exit code for CLI flows.
// The wrapped error is printed to stderr by main.
type ExitError struct {
	Code int
	Err  error
}

// Error returns the wrapped error message to satisfy error.
func (e *ExitError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}
