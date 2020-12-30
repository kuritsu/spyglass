package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
)

// API is the API object
type API struct {
	db storage.Provider
}

// Create is the API host
func Create(db storage.Provider) *API {
	result := API{db}
	return &result
}

// Serve the API
func (api *API) Serve() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	monitors := MonitorController{}
	targets := TargetController{}
	monitors.Init(api.db)
	targets.Init(api.db)

	r.GET("/monitors/:id", monitors.Get)
	r.GET("/monitors", monitors.GetAll)
	r.POST("/monitors", monitors.Post)
	r.GET("/targets/:id", targets.Get)
	r.GET("/targets", targets.GetAll)
	r.PATCH("/targets/:id", targets.Patch)
	r.POST("/targets", targets.Post)

	r.Run()
}
