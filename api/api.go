package api

import "github.com/gin-gonic/gin"

// Serve is the API host
func Serve() *gin.Engine {
	r := gin.Default()
	targets := Targets{}
	monitors := Monitors{}
	r.GET("/targets", targets.Get)
	r.GET("/monitors", monitors.Get)
	r.POST("/monitors", monitors.Post)
	return r
}
