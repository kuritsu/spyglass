package client

import "github.com/kuritsu/spyglass/api/types"

// APICaller interface
type APICaller interface {
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(*types.Target) error
}
