package main

import (
	"database/sql"
	"fmt"
	"log"

	"git.sr.ht/~arielcostas/new.vigo360.es/db"
	"github.com/joho/godotenv"
)

var (
	version string
	con     *sql.DB
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Resp struct {
	Now string
}

func main() {
	fmt.Printf("Starting vigo360 version %s\n", version)

	con = db.CreateCon()
	rows, err := con.Query("SELECT NOW()")

	if err != nil {
		fmt.Println("error querying " + err.Error())
	}

	for rows.Next() {
		var r Resp
		rows.Scan(&r.Now)
		println("row 0 => " + r.Now)
	}
}
