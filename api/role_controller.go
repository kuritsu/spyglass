package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

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

type UserRolesModifyFunc func([]string, string) []string

type ModifyUserRolesArgs struct {
	c           *gin.Context
	currentUser *types.User
	role        string
	modifyFunc  UserRolesModifyFunc
	users       []string
}

// Init -ialize the controller
func (t *RoleController) Init(db storage.Provider, log *logr.Logger) {
	t.db = db
	t.log = log
}

// GetAll roles, paginated
func (t *RoleController) GetAll(c *gin.Context) {
	t.db.Init()
	defer t.db.Free()
	pageSizeString := c.DefaultQuery("pageSize", "100")
	pageIndexString := c.DefaultQuery("pageIndex", "0")
	pageSize, err := strconv.ParseInt(pageSizeString, 10, 64)
	if err != nil || pageSize > 100 || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page size.",
		})
		return
	}
	pageIndex, err := strconv.ParseInt(pageIndexString, 10, 64)
	if err != nil || pageIndex < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page index.",
		})
		return
	}
	roles, err := t.db.GetAllRoles(pageSize, pageIndex)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	user := GetCurrentUser(c)
	result := make([]*types.Role, 0, len(roles))
	for _, r := range roles {
		if !CheckPermissions(user, r.Readers) {
			continue
		}
		result = append(result, r)
	}
	c.JSON(http.StatusOK, []*types.Role(result))
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
	user := GetCurrentUser(c)
	role.Owners = EnsurePermissions(role.Owners, user.Email)
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

func (t *RoleController) Update(c *gin.Context) {
	var roleUpdateRequest types.RoleUpdateRequest
	if er := c.ShouldBind(&roleUpdateRequest); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	roleId := c.Param("id")
	role, err := t.db.GetRole(roleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := GetCurrentUser(c)
	if !CheckPermissions(user, role.Writers) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Not enough permissions to change role.",
		})
		return
	}
	if len(roleUpdateRequest.UsersAdd) > 0 {
		t.modifyUserRoles(ModifyUserRolesArgs{c, user, roleId, t.appendRole, roleUpdateRequest.UsersAdd})
	}
	if len(roleUpdateRequest.UsersRemove) > 0 {
		t.modifyUserRoles(ModifyUserRolesArgs{c, user, roleId, t.removeRole, roleUpdateRequest.UsersRemove})
	}
	c.JSON(http.StatusOK, gin.H{
		"updateAt": time.Now(),
	})
}

func (t *RoleController) appendRole(roles []string, role string) []string {
	if slices.Contains(roles, role) {
		return roles
	}
	return append(roles, role)
}

func (t *RoleController) removeRole(roles []string, role string) []string {
	if !slices.Contains(roles, role) {
		return roles
	}
	return slices.DeleteFunc(roles, func(item string) bool { return item == role })
}

func (t *RoleController) modifyUserRoles(args ModifyUserRolesArgs) {
	for _, u := range args.users {
		tempUser, err := t.db.GetUser(u)
		if err != nil {
			args.c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if !CheckPermissions(args.currentUser, tempUser.Writers) {
			args.c.JSON(http.StatusForbidden, gin.H{
				"message": fmt.Sprintf("Not enough permissions to change user %v.", u),
			})
			return
		}
		tempUser.Roles = args.modifyFunc(tempUser.Roles, args.role)
		err = t.db.UpdateUser(tempUser, "", "")
		if err != nil {
			args.c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
}
