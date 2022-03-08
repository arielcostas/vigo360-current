package main

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

//go:embed includes
var includes embed.FS

// Map each file to its etag
var etags map[string]string = make(map[string]string)

func includesHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	inm := r.Header.Get("If-None-Match")
	if etags[file] == inm && inm != "" {
		w.WriteHeader(304)
		return
	}

	ext := regexp.MustCompile(`\.[A-Za-z]+$`).FindString(file)
	bytes, err := includes.ReadFile("includes/" + file)

	if err != nil {
		fmt.Printf("error serving file: " + err.Error() + "\n")
		http.NotFound(w, r)
	}

	// ETag for file not calculated
	if _, ok := etags[file]; !ok {
		etags[file] = GenerateEtag(bytes)
	}

	w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	w.Header().Add("Cache-Control", "max-age=2592000")
	w.Header().Add("ETag", etags[file])
	w.Write(bytes)
}

func GenerateEtag(body []byte) string {
	hasher := sha1.New()
	hasher.Write(body)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func initIncludesRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(`/includes/{file:[\w|\.|\-|\_|ñ]+}`, includesHandler).Methods("GET")
	return router
}
