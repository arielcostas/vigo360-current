package main

import (
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Information("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
