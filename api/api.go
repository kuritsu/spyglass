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
	monitors := MonitorController{}
	targets := TargetController{}
	monitors.Init(api.db, api.log)
	targets.Init(api.db, api.log)

	r.GET("/monitors/:id", monitors.Get)
	r.GET("/monitors", monitors.GetAll)
	r.POST("/monitors", monitors.Post)
	r.PUT("/monitors/:id", monitors.Put)
	r.GET("/targets/:id", targets.Get)
	r.GET("/targets", targets.GetAll)
	r.PATCH("/targets/:id", targets.Patch)
	r.POST("/targets", targets.Post)
	r.PUT("/targets/:id", targets.Put)

	return r
}
