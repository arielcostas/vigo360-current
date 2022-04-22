package admin

import "regexp"

var postIdRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\_]{3,40}$`)

func ValidatePostId(id string) bool {
	return postIdRegexp.MatchString(id)
}
