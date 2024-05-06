package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// TargetController actions
type TargetController struct {
	db  storage.Provider
	log *logr.Logger
}

// Init -ialize the controller
func (t *TargetController) Init(db storage.Provider, log *logr.Logger) {
	t.db = db
	t.log = log
}

// Get a target by its ID
func (t *TargetController) Get(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid parameter: id.",
		})
		return
	}
	includeChildrenParam := c.Query("includeChildren")
	includeChildren := false
	if includeChildrenParam != "" {
		var ok error
		includeChildren, ok = strconv.ParseBool(includeChildrenParam)
		if ok != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Invalid parameter: includeChildren.",
			})
			return
		}
	}
	t.db.Init()
	defer t.db.Free()
	target, err := t.db.GetTargetByID(id, includeChildren)
	switch {
	case err != nil:
		t.log.Error(err)
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
	userValue, _ := c.Get("user")
	user := userValue.(*types.User)
	if target.Readers != nil && len(target.Readers) > 1 && !CommonElems(target.Readers, user.Roles) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Target not found.",
		})
	}
	if target.Children != nil {
		tempChildren := make([]types.TargetRef, 0, len(target.Children))
		for _, c := range target.Children {
			if c.Readers != nil && len(c.Readers) > 1 && !CommonElems(c.Readers, user.Roles) {
				continue
			}
			tempChildren = append(tempChildren, c)
		}
		target.Children = tempChildren
	}
	c.JSON(http.StatusOK, target)
}

// GetAll targets, paginated and filtered
func (t *TargetController) GetAll(c *gin.Context) {
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
	contains := c.Query("contains")
	if contains != "" && !IsValidIDFragment(contains) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid contains expression.",
		})
		return
	}
	targets, err := t.db.GetAllTargets(pageSize, pageIndex, contains)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	userValue, _ := c.Get("user")
	user := userValue.(*types.User)
	result := make([]*types.Target, 0, len(targets))
	for _, t := range targets {
		if t.Readers != nil && len(t.Readers) > 0 && !CommonElems(t.Readers, user.Roles) {
			continue
		}
		result = append(result, t)
	}
	c.JSON(http.StatusOK, []*types.Target(result))
}

// Post a new target
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
	_, perr := t.parentMissing(target.ID)
	if perr != nil {
		c.JSON(perr.Code, gin.H{
			"message": perr.Error(),
		})
		return
	}
	if perr := t.monitorMissing(target.Monitor); perr != nil {
		c.JSON(perr.Code, gin.H{
			"message": perr.Error(),
		})
		return
	}
	t.log.Debug("Inserting in DB...")
	_, err := t.db.InsertTarget(&target)
	if err != nil {
		t.log.Error(err)
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
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

// Patch the target status
func (t *TargetController) Patch(c *gin.Context) {
	id := c.Query("id")
	var targetPatch types.TargetPatch
	er := c.ShouldBind(&targetPatch)
	switch {
	case er != nil:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	case targetPatch.Status < 0 || targetPatch.Status > 100:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid status.",
		})
		return
	case !IsValidID(id):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid target ID.",
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	target, err := t.db.GetTargetByID(id, false)
	switch {
	case err != nil:
		t.log.Error(err.Error())
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
	newTarget, err := t.db.UpdateTargetStatus(target, &targetPatch)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error. Try again.",
		})
		return
	}
	c.JSON(http.StatusOK, newTarget)
}

// Put an existing target.
func (t *TargetController) Put(c *gin.Context) {
	id := c.Param("id")
	forceStatusUpdate := c.Query("forceStatusUpdate") == "true"
	var targetObj types.Target
	er := c.ShouldBind(&targetObj)
	switch {
	case er != nil:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	case targetObj.Status < 0 || targetObj.Status > 100:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid status.",
		})
		return
	case !IsValidID(id):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid target ID.",
		})
		return
	}
	t.db.Init()
	defer t.db.Free()
	target, err := t.db.GetTargetByID(id, false)
	switch {
	case err != nil:
		t.log.Error(err)
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
	newTarget, err := t.db.UpdateTarget(target, &targetObj, forceStatusUpdate)
	if err != nil {
		t.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error. Try again.",
		})
		return
	}
	c.JSON(http.StatusOK, newTarget)
}

func (t *TargetController) parentMissing(childID string) (*types.Target, *ErrorWithCode) {
	parentID := types.GetTargetParentByID(childID)
	if parentID == "" {
		return nil, nil
	}

	parent, err := t.db.GetTargetByID(parentID, false)
	if err != nil {
		t.log.Error(err)
		return nil, &ErrorWithCode{Message: "Invalid operation. Try again later", Code: http.StatusInternalServerError}
	}
	if parent == nil {
		msg := fmt.Sprintf("Target parent does not exist. (%s)", parentID)
		t.log.Println(msg)
		return nil, &ErrorWithCode{Message: msg, Code: http.StatusBadRequest}
	}
	return parent, nil
}

func (t *TargetController) monitorMissing(monitorRef *types.MonitorRef) *ErrorWithCode {
	if monitorRef == nil {
		return nil
	}
	monitor, err := t.db.GetMonitorByID(monitorRef.MonitorID)
	if err != nil {
		t.log.Error(err)
		return &ErrorWithCode{Message: "Invalid operation. Try again later", Code: http.StatusInternalServerError}
	}
	if monitor == nil {
		msg := fmt.Sprintf("Monitor does not exist. (%s)", monitorRef.MonitorID)
		t.log.Println(msg)
		return &ErrorWithCode{Message: msg, Code: http.StatusBadRequest}
	}
	return nil
}
