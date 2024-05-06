package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	logr "github.com/sirupsen/logrus"
)

func abortWithMessage(c *gin.Context, msg string, db storage.Provider) {
	if db != nil {
		db.Free()
	}
	c.JSON(http.StatusForbidden, gin.H{
		"message": msg,
	})
	c.Abort()
}

func AuthMiddleware(db storage.Provider, log *logr.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Request.Header.Get("Authorization")
		before, after, found := strings.Cut(h, ":")
		if h == "" || !found {
			abortWithMessage(c, "Unauthorized", nil)
			return
		}
		log.Debug("Searching for user ", before)
		db.Init()
		user, err := db.GetUser(before)
		if err != nil {
			log.Error(err)
			abortWithMessage(c, "Invalid user", db)
			return
		}
		c.Set("user", user)
		log.Debug("Validating token for ", before)
		err = db.ValidateToken(before, after)
		log.Debug("Finished validation.")
		db.Free()
		if err != nil {
			log.Error(err)
			abortWithMessage(c, "Unauthorized", nil)
			return
		}
		c.Next()
	}
}
