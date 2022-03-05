package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/admin"
	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/public"
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

	common.DatabaseInit()

	http.Handle("/admin/", admin.InitRouter())
	http.Handle("/includes/", initIncludesRouter())
	http.Handle("/", public.InitRouter())

	log.Fatal(http.ListenAndServe(PORT, nil))
}
