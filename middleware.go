package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s", time.Now().Format("15:04:05"), r.Method, r.RequestURI)
		if !strings.HasSuffix(r.Referer(), r.RequestURI) {
			fmt.Printf(" - Ref: %s\n", r.Referer())
		} else {
			fmt.Print("\n")
		}
		next.ServeHTTP(w, r)
	})
}
