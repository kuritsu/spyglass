package api

import (
	"net/http"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// UserController actions
type RoleController struct {
	db  storage.Provider
	log *logr.Logger
}

// Init -ialize the controller
func (t *RoleController) Init(db storage.Provider, log *logr.Logger) {
	t.db = db
	t.log = log
}

func (t *RoleController) Add(c *gin.Context) {
	var role types.Role
	if er := c.ShouldBind(&role); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	userValue, _ := c.Get("user")
	user := userValue.(*types.User)
	role.Owners = EnsurePermissions(role.Owners, user.Email)
	role.Readers = EnsurePermissions(role.Readers, user.Email)
	role.Writers = EnsurePermissions(role.Writers, user.Email)
	err := t.db.InsertRole(&role, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"updateAt": role.UpdatedAt,
	})
}
