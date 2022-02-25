package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/routes"
	"github.com/joho/godotenv"
)

var (
	version string
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Printf("Starting vigo360 version %s\n", version)
	var PORT string = ":" + os.Getenv("PORT")

	fmt.Println("Starting web server on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, routes.InitRouter()))
}
