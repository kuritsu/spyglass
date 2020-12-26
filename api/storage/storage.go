package storage

// CreateProviderFromConf returns a provider according to the configuration.
func CreateProviderFromConf() Provider {
	var result Provider
	result = new(MongoDB)
	return result
}
