package main

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"mime"
	"net/http"
	"regexp"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

//go:embed includes
var includes embed.FS

// Map each file to its etag
var etags map[string]string = make(map[string]string)

func pregenerateEtags() {
	entries, _ := includes.ReadDir("includes")
	for _, e := range entries {
		filename := e.Name()
		bytes, err := includes.ReadFile("includes/" + filename)
		if err != nil { // Log the error and skip to next one
			logger.Error("error pregenerating etag for %s: %s", filename, err.Error())
			continue
		}
		etags[filename] = GenerateEtag(bytes)
	}
	logger.Information("pregenerated etags for %d includes", len(entries))
}

func includesHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	inm := r.Header.Get("If-None-Match")

	// The client has the current version of the file, so we reply with 304 (Not Modified) and go on with our lives
	if etags[file] == inm && inm != "" {
		w.WriteHeader(304)
		return
	}

	// We need the extension to return the correct MIME, needed by browsers
	ext := regexp.MustCompile(`\.[A-Za-z]+$`).FindString(file)
	bytes, err := includes.ReadFile("includes/" + file)

	if err != nil {
		logger.Error("[includes]: error serving file %s: %s", file, err.Error())
		w.WriteHeader(404)
		return
	}

	// Long-term cache, to reduce server load and bandwidth consumption
	w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("ETag", etags[file])
	_, err = w.Write(bytes)
	if err != nil {
		logger.Error("[includes] error writing file %s: %", file, err.Error())
		return
	}
}

// Receives the contents of a file and returns a base64-encoded SHA-1 of the file, suitable for the ETag HTTP header
func GenerateEtag(body []byte) string {
	hasher := sha1.New()
	_, err := hasher.Write(body)
	if err != nil {
		logger.Error("[load] error generating body hash: %s", err.Error())
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Initialize the module, pregenerating the ETags and the router
func initIncludesRouter() *mux.Router {
	pregenerateEtags()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(`/includes/{file:[\w|\.|\-|\_|ñ]+}`, includesHandler).Methods("GET")
	return router
}
