package api

import (
	"github.com/gin-gonic/gin"
	logr "github.com/sirupsen/logrus"

	"github.com/kuritsu/spyglass/api/storage"
)

// API is the API object
type API struct {
	db  storage.Provider
	log *logr.Logger
}

// Create is the API host
func Create(db storage.Provider, log *logr.Logger) *API {
	result := API{db, log}
	return &result
}

// Serve the API
func (api *API) Serve() *gin.Engine {
	if api.log.Level == logr.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	gin.DisableConsoleColor()
	r.Use(CORSMiddleware())
	monitors := MonitorController{}
	targets := TargetController{}
	users := UserController{}
	monitors.Init(api.db, api.log)
	targets.Init(api.db, api.log)
	users.Init(api.db, api.log)

	authMid := AuthMiddleware(api.db, api.log)

	r.GET("/monitors/:id", authMid, monitors.Get)
	r.GET("/monitors", authMid, monitors.GetAll)
	r.POST("/monitors", authMid, monitors.Post)
	r.PUT("/monitors/:id", authMid, monitors.Put)
	r.GET("/target", authMid, targets.Get)
	r.GET("/targets", authMid, targets.GetAll)
	r.PATCH("/target", authMid, targets.Patch)
	r.POST("/targets", authMid, targets.Post)
	r.PUT("/target", authMid, targets.Put)
	r.POST("/login", users.Login)
	r.POST("/register", users.Register)

	return r
}
