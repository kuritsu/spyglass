package client

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
)

// APICaller interface
type APICaller interface {
	CreateUserToken(email string, expiration time.Time) (string, error)
	DeleteTarget(string) (int, error)
	GetTargetByID(id string, includeChildren bool) (types.TargetRef, error)
	Init(string)
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error
	InsertRole(*types.Role) error
	ListMonitors(int, int) ([]*types.Monitor, error)
	ListRoles(pageIndex, pageSize int) ([]*types.Role, error)
	ListTargets(string, int, int) ([]*types.Target, error)
	ListUsers(pageIndex, pageSize int) ([]*types.User, error)
	Login(string, string) (string, error)
	UpdateRole(role string, usersAdd, usersRemove []string) error
	UpdateTargetStatus(string, int, string) error
}
