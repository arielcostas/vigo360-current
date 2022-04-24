/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import "regexp"

var postIdRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\_]{3,40}$`)

func ValidatePostId(id string) bool {
	return postIdRegexp.MatchString(id)
}
