package storage

// CreateProviderFromConf returns a provider according to the configuration.
func CreateProviderFromConf() Provider {
	result := MongoDB{}
	return &result
}
