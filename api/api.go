package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
)

// Serve is the API host
func Serve(db storage.Provider) *gin.Engine {
	r := gin.Default()

	monitors := MonitorController{}
	targets := TargetController{}
	monitors.Init(db)
	targets.Init(db)

	r.GET("/monitors/:id", monitors.Get)
	r.GET("/monitors", monitors.GetAll)
	r.POST("/monitors", monitors.Post)
	r.GET("/targets/:id", targets.Get)
	r.GET("/targets", targets.GetAll)
	r.POST("/targets", targets.Post)
	return r
}
