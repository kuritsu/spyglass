package api

import (
	"net/http"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// UserController actions
type UserController struct {
	db  storage.Provider
	log *logr.Logger
}

// Init -ialize the controller
func (t *UserController) Init(db storage.Provider, log *logr.Logger) {
	t.db = db
	t.log = log
}

// Get a target by its ID
func (t *UserController) Login(c *gin.Context) {
	var creds types.AuthRequest
	if er := c.ShouldBind(&creds); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	user, err := t.db.Login(creds.Email, creds.Password)
	if err != nil {
		t.log.Error(err)
		if err.Error() == "InvalidCredentials" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid credentials.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	t.log.Info("User ", creds.Email, " authenticated.")
	token, err := t.db.CreateUserToken(user)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, token)
}
