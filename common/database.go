package common

import (
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Database *sqlx.DB

func DatabaseInit() {
	var dsn string = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_BASE")
	var err error
	Database, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Critical("error connecting to mysql: " + err.Error())
	}

	logger.Information("database connection established")

	err = Database.Ping()
	if err != nil {
		logger.Critical("couldn't ping database: " + err.Error())
	}

	Database.Exec("SET lc_time_names = 'es_ES';")
	logger.Information("database configured")
}
