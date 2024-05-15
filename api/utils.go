package api

import (
	"encoding/json"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/types"
	"gopkg.in/yaml.v2"
)

func CommonElems(first, second []string) bool {
	set := make(map[string]bool)
	for _, e := range first {
		set[e] = true
	}
	for _, e := range second {
		_, ok := set[e]
		if ok {
			return true
		}
	}
	return false
}

func NewObjectFromFile[T any](fileName string) (*T, error) {
	lowerFileName := strings.ToLower(fileName)
	var result T
	var unmarshalFunc func(in []byte, out any) (err error)
	switch {
	case strings.HasSuffix(lowerFileName, ".json"):
		unmarshalFunc = json.Unmarshal
	case strings.HasSuffix(lowerFileName, ".yaml") || strings.HasSuffix(lowerFileName, ".yml"):
		unmarshalFunc = yaml.Unmarshal
	default:
		return &result, os.ErrInvalid
	}
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return &result, err
	}
	err = unmarshalFunc(raw, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func CheckPermissions(user *types.User, targetPermissions []string) bool {
	if len(targetPermissions) > 0 &&
		!slices.Contains(user.Roles, "admins") &&
		!slices.Contains(targetPermissions, user.Email) &&
		!CommonElems(targetPermissions, user.Roles) {
		return false
	}
	return true
}

func EnsurePermissions(perms []string, user string) []string {
	if len(perms) == 0 {
		return []string{user}
	}
	return perms
}

func GetCurrentUser(c *gin.Context) *types.User {
	userValue, _ := c.Get("user")
	user := userValue.(*types.User)
	return user
}
