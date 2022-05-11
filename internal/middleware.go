/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thanhpk/randstr"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rid = randstr.String(10)
		fmt.Printf("<6>[%s] [%s] %s %s\n", r.Header.Get("X-Forwarded-For"), rid, r.Method, r.URL.Path)
		newContext := context.WithValue(r.Context(), ridContextKey("rid"), rid)
		r = r.WithContext(newContext)

		next.ServeHTTP(w, r)
	})
}
