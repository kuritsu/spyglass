package storage

import "github.com/kuritsu/spyglass/api/types"

// Provider for storage
type Provider interface {
	GetMonitor(string) *types.Monitor
}
