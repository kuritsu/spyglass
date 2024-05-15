package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	token, err := t.db.CreateUserToken(user, time.Now().UTC().Add(time.Hour*24))
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
	token, err := t.db.CreateUserToken(user, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (t *UserController) CreateToken(c *gin.Context) {
	var req types.UserTokenRequest
	if er := c.ShouldBind(&req); er != nil || req.Expiration.After(time.Now().AddDate(1, 0, 0)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	userId := c.Param("id")
	user, err := t.db.GetUser(userId)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	currentUser := GetCurrentUser(c)
	if !CheckPermissions(currentUser, user.Writers) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("Not enough permissions on user %s.", userId),
		})
		return
	}
	token, err := t.db.CreateUserToken(user, req.Expiration)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, token)
}

// GetAll users, paginated
func (t *UserController) GetAll(c *gin.Context) {
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
	users, err := t.db.GetAllUsers(pageSize, pageIndex)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	user := GetCurrentUser(c)
	result := make([]*types.User, 0, len(users))
	for _, r := range users {
		if !CheckPermissions(user, r.Readers) {
			continue
		}
		r.PassHash = ""
		r.FirstHash = ""
		result = append(result, r)
	}
	c.JSON(http.StatusOK, []*types.User(result))
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
	user := GetCurrentUser(c)
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
