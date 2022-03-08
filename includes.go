package main

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

//go:embed includes
var includes embed.FS

// Map each file to its etag
var etags map[string]string = make(map[string]string)

func pregenerateEtags() {
	entries, _ := includes.ReadDir("includes")
	for _, e := range entries {
		bytes, err := includes.ReadFile("includes/" + e.Name())
		if err != nil {
			log.Fatalf("error reading " + e.Name())
		}
		etags[e.Name()] = GenerateEtag(bytes)
	}
	fmt.Printf("Pregenerated etags for %d includes\n", len(entries))
}

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
		println("calculating etag for file " + file)
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
	pregenerateEtags()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(`/includes/{file:[\w|\.|\-|\_|ñ]+}`, includesHandler).Methods("GET")
	return router
}
