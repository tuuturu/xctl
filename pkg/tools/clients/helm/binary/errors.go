package binary

import "regexp"

var reUnreachableErr = regexp.MustCompile(`.*EOF: Kubernetes cluster unreachable.*`)

func isUnreachable(err error) bool {
	return reUnreachableErr.MatchString(err.Error())
}

var reTimedOutErr = regexp.MustCompile(`.*connection timed out.*`)

func isConnectionTimedOut(err error) bool {
	return reTimedOutErr.MatchString(err.Error())
}

var reAlreadyExistsErr = regexp.MustCompile(`INSTALLATION FAILED: cannot re-use a name that is still in use\s`)

func isAlreadyExists(err error) bool {
	return reAlreadyExistsErr.MatchString(err.Error())
}
