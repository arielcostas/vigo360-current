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

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/admin"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/public"
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
			var rid = randstr.String(30)
			ctx = context.WithValue(ctx, "rid", rid)
			r = r.WithContext(ctx)

			logger.Information("[%s] %s - %s %s", rid, r.Header["X-Forwarded-For"], r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
	return r
}

func main() {
	logger.Information("starting Vigo360 version " + version)
	var PORT string = ":" + os.Getenv("PORT")

	logger.Information("starting web server on %s", PORT)

	var db = database.GetDB()

	// Automatically revoke sessions every 6 hours
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(5).Minutes().Do(func() {
		res, err := db.DB.Exec(`UPDATE sesiones SET revocada = 0 WHERE iniciada < DATE_SUB(NOW(), INTERVAL 6 HOUR);`)
		if err != nil {
			logger.Critical("error cleaning old sessions: %s", err.Error())
		}
		ra, err := res.RowsAffected()
		if err != nil {
			logger.Critical("error displaying cleaned sessions: %s", err.Error())
		}

		if ra > 0 {
			logger.Information("automatically revoked %d session(s)", ra)
		}
	})
	if err != nil {
		logger.Critical("error running session deleter: %s", err.Error())
	}
	s.StartAsync()

	http.Handle("/admin/", mw(admin.InitRouter()))
	http.Handle("/", mw(public.InitRouter()))

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Critical("error with HTTP server: %s", err.Error())
	}
}
