package client

import "github.com/kuritsu/spyglass/api/types"

// APICaller interface
type APICaller interface {
	Init(string)
	Login(string, string) (string, error)
	ListTargets(string, int, int) ([]*types.Target, error)
	GetTargetByID(id string, includeChildren bool) (types.TargetRef, error)
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error
	UpdateTargetStatus(string, int, string) error
}
