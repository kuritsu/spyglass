package api

import "regexp"

var idRegexp *regexp.Regexp
var idFragRegexp *regexp.Regexp

func init() {
	idRegexp, _ = regexp.Compile(`^[\w\d\-_]+(\.[\w\d\-_]+)*$`)
	idFragRegexp, _ = regexp.Compile(`^[\w\d\-_\.]+$`)
}

// IsValidID determines if an ID is valid
func IsValidID(id string) bool {
	return idRegexp.MatchString(id)
}

// IsValidIDFragment determines if a string is a valid ID fragment
func IsValidIDFragment(id string) bool {
	return idFragRegexp.MatchString(id)
}
