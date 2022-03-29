/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package common

import (
	"net/http"
	"strings"
)

func Redirect(w http.ResponseWriter, r *http.Request, from string, to string) {
	http.Redirect(w, r,
		strings.ReplaceAll(r.URL.String(), from, to),
		http.StatusMovedPermanently)
}
