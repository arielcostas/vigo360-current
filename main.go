package main

import (
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/admin"
	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"git.sr.ht/~arielcostas/new.vigo360.es/public"
	"github.com/gorilla/mux"
)

var (
	version string
)

func wrapMiddleware(r *mux.Router) *mux.Router {
	r.Use(LogMiddleware)
	return r
}

func main() {
	logger.Information("starting Vigo360 version " + version)
	var PORT string = ":" + os.Getenv("PORT")

	logger.Information("starting web server on %s", PORT)

	common.DatabaseInit()

	http.Handle("/admin/", wrapMiddleware(admin.InitRouter()))
	http.Handle("/includes/", wrapMiddleware(initIncludesRouter()))
	http.Handle("/", wrapMiddleware(public.InitRouter()))

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Critical("error with HTTP server: %s", err.Error())
	}
}
