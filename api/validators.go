package api

import "regexp"

var idRegexp *regexp.Regexp

func init() {
	idRegexp, _ = regexp.Compile(`^[\w\d\-_]+(\.[\w\d\-_]+)*$`)
}

// IsValidID determines if an ID is valid
func IsValidID(id string) bool {
	return idRegexp.MatchString(id)
}
