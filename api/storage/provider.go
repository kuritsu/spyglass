package storage

import "github.com/kuritsu/spyglass/api/types"

// Provider for storage
type Provider interface {
	Init()
	Free()

	GetMonitorByID(string) (*types.Monitor, error)
	GetTargetByID(string) (*types.Target, error)
	InsertMonitor(*types.Monitor) (*types.Monitor, error)
	InsertTarget(*types.Target) (*types.Target, error)
}
