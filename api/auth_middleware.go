package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	logr "github.com/sirupsen/logrus"
)

func AuthMiddleware(db storage.Provider, log *logr.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Request.Header.Get("Authorization")
		if h == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		before, after, found := strings.Cut(h, ":")
		if !found {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		db.Init()
		defer db.Free()
		log.Debug("Searching for user ", before)
		user, err := db.GetUser(before)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid user",
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		log.Debug("Validating token for ", before)
		err = db.ValidateToken(before, after)
		log.Debug("Finished validation.")
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
