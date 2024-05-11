package client

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
)

// APICaller interface
type APICaller interface {
	Init(string)
	Login(string, string) (string, error)
	ListTargets(string, int, int) ([]*types.Target, error)
	GetTargetByID(id string, includeChildren bool) (types.TargetRef, error)
	InsertOrUpdateMonitor(*types.Monitor) error
	InsertOrUpdateTarget(target *types.Target, forceStatusUpdate bool) error
	UpdateTargetStatus(string, int, string) error
	InsertRole(*types.Role) error
	UpdateRole(role string, usersAdd, usersRemove []string) error
	ListRoles(pageIndex, pageSize int) ([]*types.Role, error)
	ListUsers(pageIndex, pageSize int) ([]*types.User, error)
	CreateUserToken(email string, expiration time.Time) (string, error)
}
