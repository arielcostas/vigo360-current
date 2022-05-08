/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/admin"
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
	if err := checkEnv(); err != nil {
		logger.Critical("error validando entorno: %s\n", err.Error())
	}

	if err := run(); err != nil {
		logger.Critical("%s\n", err)
	}
}

func run() error {
	logger.Information("starting Vigo360 version " + version)
	var PORT string = ":" + os.Getenv("PORT")

	logger.Information("starting web server on %s", PORT)

	http.Handle("/admin/", mw(admin.InitRouter()))
	http.Handle("/", mw(public.InitRouter()))

	var err = http.ListenAndServe(PORT, nil)
	return err
}

func checkEnv() error {
	if val, is := os.LookupEnv("PORT"); !is || val == "" {
		return fmt.Errorf("es necesario especificar PORT")
	} else {
		i, e := strconv.Atoi(val)
		if e != nil {
			return fmt.Errorf("PORT tiene que ser un número")
		}
		if i < 0 || i > 65535 {
			return fmt.Errorf("PORT debe ser un puerto TCP válido")
		}
	}

	if val, is := os.LookupEnv("UPLOAD_PATH"); !is || val == "" {
		return fmt.Errorf("es necesario especificar UPLOAD_PATH")
	} else {
		info, err := os.Stat(val)
		if err != nil {
			return fmt.Errorf("error comprobando validez de UPLOAD_PATH: %s", err.Error())
		}
		if !info.IsDir() {
			return fmt.Errorf("UPLOAD_PATH tiene que ser un directorio: %s", err.Error())
		}
		err = os.WriteFile(val+"/.test", []byte{0x00}, os.ModePerm)
		if err != nil {
			return fmt.Errorf("no se puede escribir a UPLOAD_PATH: %s", err.Error())
		}
		os.Remove(val + "/.test")
	}

	if val, is := os.LookupEnv("DOMAIN"); !is || val == "" {
		return fmt.Errorf("es necesario especificar DOMAIN")
	}

	return nil
}
