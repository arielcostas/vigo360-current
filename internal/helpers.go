package internal

import "os"

func baseUrl() string {
	return os.Getenv("DOMAIN")
}

func fullCanonica(path string) string {
	return baseUrl() + path
}

func getMinimo(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
