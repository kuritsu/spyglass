package storage

import "github.com/kuritsu/spyglass/api/types"

// Provider for storage
type Provider interface {
	Initialize()
	GetMonitorByID(string) *types.Monitor
	GetTargetByID(string) *types.Target
}
