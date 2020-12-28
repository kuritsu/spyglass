package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// TargetController actions
type TargetController struct {
	db storage.Provider
}

// Init -ialize the controller
func (t *TargetController) Init(db storage.Provider) {
	t.db = db
}

// Get a target by its ID
func (t *TargetController) Get(c *gin.Context) {
	id := c.Param("id")
	t.db.Init()
	defer t.db.Free()
	target, err := t.db.GetTargetByID(id)
	switch {
	case err != nil:
		log.Println("ERROR: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error. Try again.",
		})
		return
	case target == nil:
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Target not found.",
		})
		return
	}
	c.JSON(http.StatusOK, target)
}

// Post a new monitor
func (t *TargetController) Post(c *gin.Context) {
	var target types.Target
	if er := c.ShouldBind(&target); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	if !IsValidID(target.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid target ID.",
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	perr := t.parentMissing(target.ID)
	if perr != nil {
		c.JSON(perr.Code, gin.H{
			"message": perr.Error(),
		})
		return
	}
	_, err := t.db.InsertTarget(&target)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "Duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Duplicate target ID.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusCreated, target)
}

func (t *TargetController) parentMissing(childID string) *ErrorWithCode {
	parentID := types.GetTargetParentByID(childID)
	if parentID == "" {
		return nil
	}

	parent, err := t.db.GetTargetByID(parentID)
	if err != nil {
		log.Println(err.Error())
		return &ErrorWithCode{Message: "Invalid operation. Try again later", Code: http.StatusInternalServerError}
	}
	if parent == nil {
		msg := fmt.Sprintf("Target parent does not exist. (%s)", parentID)
		log.Println(msg)
		return &ErrorWithCode{Message: msg, Code: http.StatusBadRequest}
	}
	return nil
}
