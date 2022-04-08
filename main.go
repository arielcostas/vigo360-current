/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"git.sr.ht/~arielcostas/new.vigo360.es/admin"
	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"git.sr.ht/~arielcostas/new.vigo360.es/public"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
)

var (
	version string
)

func mw(r *mux.Router) *mux.Router {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
			defer cancel()

			// Generates a random RequestID to print in logs, errors and stuff
			ctx = context.WithValue(ctx, "rid", randstr.String(14))
			r = r.WithContext(ctx)

			logger.Information("%s - %s %s", r.Header["X-Forwarded-For"], r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
	return r
}

func main() {
	logger.Information("starting Vigo360 version " + version)
	var PORT string = ":" + os.Getenv("PORT")

	logger.Information("starting web server on %s", PORT)

	common.DatabaseInit()

	http.Handle("/admin/", mw(admin.InitRouter()))
	http.Handle("/includes/", mw(initIncludesRouter()))
	http.Handle("/", mw(public.InitRouter()))

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Critical("error with HTTP server: %s", err.Error())
	}
}
