// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

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
