package storage

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
)

// Provider for storage
type Provider interface {
	Init()
	Free()

	GetAllMonitors(int64, int64, string) ([]types.Monitor, error)
	GetAllTargets(int64, int64, string) ([]*types.Target, error)
	GetAllRoles(int64, int64) ([]*types.Role, error)
	GetAllUsers(int64, int64) ([]*types.User, error)
	GetMonitorByID(string) (*types.Monitor, error)
	GetTargetByID(id string, includeChildren bool) (*types.Target, error)
	InsertMonitor(*types.Monitor) (*types.Monitor, error)
	InsertTarget(*types.Target) (*types.Target, error)
	UpdateMonitor(*types.Monitor, *types.Monitor) (*types.Monitor, error)
	UpdateTargetStatus(*types.Target, *types.TargetPatch) (*types.Target, error)
	UpdateTarget(*types.Target, *types.Target, bool) (*types.Target, error)
	DeleteTarget(id string) (int, error)
	Login(string, string) (*types.User, error)
	Register(string, string) (*types.User, error)
	CreateUserToken(*types.User, time.Time) (string, error)
	ValidateToken(string, string) error
	GetUser(string) (*types.User, error)
	UpdateUser(*types.User, string, string) error
	InsertRole(*types.Role, *types.User) error
	GetRole(string) (*types.Role, error)
}
