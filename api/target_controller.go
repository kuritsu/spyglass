package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// GetAll targets, paginated and filtered
func (t *TargetController) GetAll(c *gin.Context) {
	t.db.Init()
	defer t.db.Free()
	pageSizeString := c.DefaultQuery("pageSize", "10")
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
	contains := c.Query("contains")
	if contains != "" && !IsValidIDFragment(contains) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid contains expression.",
		})
		return
	}
	targets, err := t.db.GetAllTargets(pageSize, pageIndex, contains)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusOK, []types.Target(targets))
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
	perr = t.monitorMissing(target.Monitor)
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

func (t *TargetController) monitorMissing(monitorRef *types.MonitorRef) *ErrorWithCode {
	if monitorRef == nil {
		return nil
	}
	monitor, err := t.db.GetMonitorByID(monitorRef.MonitorID)
	if err != nil {
		log.Println(err.Error())
		return &ErrorWithCode{Message: "Invalid operation. Try again later", Code: http.StatusInternalServerError}
	}
	if monitor == nil {
		msg := fmt.Sprintf("Monitor does not exist. (%s)", monitorRef.MonitorID)
		log.Println(msg)
		return &ErrorWithCode{Message: msg, Code: http.StatusBadRequest}
	}
	return nil
}
