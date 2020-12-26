package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
)

// Serve is the API host
func Serve(db storage.Provider) *gin.Engine {
	r := gin.Default()
	targets := Targets{}
	monitors := MonitorController{}
	monitors.Initialize(db)
	r.GET("/targets", targets.Get)
	r.GET("/monitors", monitors.Get)
	r.POST("/monitors", monitors.Post)
	return r
}
