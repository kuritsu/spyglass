package client

import "github.com/kuritsu/spyglass/api/types"

// APICaller interface
type APICaller interface {
	Init(string)
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(*types.Target, bool) error
}
