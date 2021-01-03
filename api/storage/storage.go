package storage

import "github.com/sirupsen/logrus"

// CreateProviderFromConf returns a provider according to the configuration.
func CreateProviderFromConf(log *logrus.Logger) Provider {
	result := MongoDB{
		Log: log,
	}
	return &result
}
