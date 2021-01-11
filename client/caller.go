package client

import "github.com/kuritsu/spyglass/api/types"

// APICaller interface
type APICaller interface {
	Init(string)
	ListTargets(string, int, int) ([]*types.Target, error)
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(*types.Target, bool) error
	UpdateTargetStatus(string, int, string) error
}
