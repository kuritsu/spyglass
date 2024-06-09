package storage

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
)

// Provider for storage
type Provider interface {
	Init()
	Seed()
	Free()

	CreateUserToken(*types.User, time.Time) (string, error)
	DeleteJob(id string) error
	DeleteScheduler(id string) error
	DeleteTarget(id string) (int, error)
	GetAllJobsFor(label string) ([]*types.Job, error)
	GetAllMonitors(int64, int64, string) ([]*types.Monitor, error)
	GetAllRoles(int64, int64) ([]*types.Role, error)
	GetAllSchedulersFor(label string) ([]*types.Scheduler, error)
	GetAllTargets(int64, int64, string) ([]*types.Target, error)
	GetAllUsers(int64, int64) ([]*types.User, error)
	GetMonitorByID(string) (*types.Monitor, error)
	GetRole(string) (*types.Role, error)
	GetTargetByID(id string, includeChildren bool) (*types.Target, error)
	GetUser(string) (*types.User, error)
	InsertJob(*types.Job) (*types.Job, error)
	InsertMonitor(*types.Monitor) (*types.Monitor, error)
	InsertRole(*types.Role, *types.User) error
	InsertScheduler(*types.Scheduler) (*types.Scheduler, error)
	InsertTarget(*types.Target) (*types.Target, error)
	Login(string, string) (*types.User, error)
	Register(string, string) (*types.User, error)
	UpdateJob(*types.Job) (*types.Job, error)
	UpdateMonitor(*types.Monitor, *types.Monitor) (*types.Monitor, error)
	UpdateScheduler(*types.Scheduler) (*types.Scheduler, error)
	UpdateTarget(*types.Target, *types.Target, bool) (*types.Target, error)
	UpdateTargetStatus(*types.Target, *types.TargetPatch) (*types.Target, error)
	UpdateUser(*types.User, string, string) error
	ValidateToken(string, string) error
}
