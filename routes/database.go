package routes

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() {
	var dsn string = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_BASE")
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connected to database")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db is not connected")
		fmt.Println(err.Error())
	}

	db.Exec("SET lc_time_names = 'es_ES';")
}
