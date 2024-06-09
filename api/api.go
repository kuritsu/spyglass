package api

import (
	"github.com/gin-gonic/gin"
	logr "github.com/sirupsen/logrus"

	"github.com/kuritsu/spyglass/api/storage"
)

// API is the API object
type API struct {
	db              storage.Provider
	log             *logr.Logger
	statusUpdateJob *StatusUpdateJob
}

// Create is the API host
func Create(db storage.Provider, log *logr.Logger, statusUpdateJob *StatusUpdateJob) *API {
	result := API{db, log, statusUpdateJob}
	return &result
}

// Serve the API
func (api *API) Serve() *gin.Engine {
	api.db.Init()
	api.db.Seed()
	api.db.Free()

	if api.log.Level == logr.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	gin.DisableConsoleColor()
	r.Use(CORSMiddleware())
	monitors := MonitorController{}
	targets := TargetController{
		statusUpdateChan: api.statusUpdateJob.StatusChan,
	}
	users := UserController{}
	roles := RoleController{}
	monitors.Init(api.db, api.log)
	targets.Init(api.db, api.log)
	users.Init(api.db, api.log)
	roles.Init(api.db, api.log)

	authMid := AuthMiddleware(api.db, api.log)

	r.DELETE("/target", authMid, targets.Delete)
	r.GET("/monitors", authMid, monitors.GetAll)
	r.GET("/monitors/:id", authMid, monitors.Get)
	r.GET("/roles", authMid, roles.GetAll)
	r.GET("/target", authMid, targets.Get)
	r.GET("/targets", authMid, targets.GetAll)
	r.GET("/users", authMid, users.GetAll)
	r.PATCH("/role/:id", authMid, roles.Update)
	r.PATCH("/target", authMid, targets.Patch)
	r.PATCH("/user/:id", authMid, users.Update)
	r.POST("/monitors", authMid, monitors.Post)
	r.POST("/roles", authMid, roles.Add)
	r.POST("/targets", authMid, targets.Post)
	r.POST("/user/login", users.Login)
	r.POST("/user/token/:id", authMid, users.CreateToken)
	r.POST("/users", users.Register)
	r.PUT("/monitors/:id", authMid, monitors.Put)
	r.PUT("/target", authMid, targets.Put)

	go api.statusUpdateJob.Run()

	return r
}
