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
