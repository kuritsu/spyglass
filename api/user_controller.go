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

func (t *UserController) Register(c *gin.Context) {
	var creds types.AuthRequest
	if er := c.ShouldBind(&creds); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	user, err := t.db.Register(creds.Email, creds.Password)
	if err != nil {
		t.log.Error(err)
		if err.Error() == "ErrorCreatingUser" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Error creating user.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	t.log.Info("User ", creds.Email, " created.")
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

func (t *UserController) Update(c *gin.Context) {
	var userUpdate types.UserUpdateRequest
	if er := c.ShouldBind(&userUpdate); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	userValue, _ := c.Get("user")
	user := userValue.(*types.User)
	if userUpdate.FullName != "" {
		user.FullName = userUpdate.FullName
	}
	err := t.db.UpdateUser(user, userUpdate.OldPassword, userUpdate.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"updateAt": user.Permissions.UpdatedAt,
	})
}
