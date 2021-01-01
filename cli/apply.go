package cli

// ApplyOptions according to arguments
type ApplyOptions struct {
	Recursive bool
}

// Apply the configuration in the given directory.
func (c *CommandLine) Apply(dir string, options ApplyOptions) {
	c.log.Debug("Executing apply.")
}
