package linode

import "regexp"

var reNotAuthenticated = regexp.MustCompile(`.*\[401] Invalid Token$`)

func errorIsNotAuthenticated(err error) bool {
	return reNotAuthenticated.MatchString(err.Error())
}
