package internal

import (
	"html/template"
	"os"
)

func baseUrl() template.URL {
	return template.URL(os.Getenv("DOMAIN"))
}

func fullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func getMinimo(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
