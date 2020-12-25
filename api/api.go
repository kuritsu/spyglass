package api

import "github.com/gin-gonic/gin"

// Serve is the API host
func Serve() {
	r := gin.Default()
	targets := Targets{}
	monitors := Monitors{}
	r.GET("/targets", targets.Get)
	r.GET("/monitors", monitors.Get)
	r.Run()
}
