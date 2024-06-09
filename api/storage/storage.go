package storage

import (
	"github.com/kuritsu/spyglass/api/storage/mongodb"
	"github.com/sirupsen/logrus"
)

// CreateProviderFromConf returns a provider according to the configuration.
func CreateProviderFromConf(log *logrus.Logger) Provider {
	result := mongodb.MongoDB{
		Log: log,
	}
	return &result
}
