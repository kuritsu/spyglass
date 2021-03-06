package runner

import "github.com/gin-gonic/gin"

// Gin runner.
type Gin struct {
	Engine  *gin.Engine
	Address string
}

// Run engine.
func (g *Gin) Run() error {
	return g.Engine.Run(g.Address)
}
