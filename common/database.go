package common

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Database *sqlx.DB

func DatabaseInit() {
	var dsn string = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_BASE")
	var err error
	Database, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connected to database")
	}

	err = Database.Ping()
	if err != nil {
		fmt.Println("db is not connected")
		fmt.Println(err.Error())
	}

	Database.Exec("SET lc_time_names = 'es_ES';")
}
