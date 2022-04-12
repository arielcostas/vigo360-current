/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package database

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/logger"
)

var db *sqlx.DB

func start() {
	var dsn string = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_BASE")
	var err error
	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Critical("error connecting to mysql: %s", err.Error())
	}

	logger.Information("database connection established")

	err = conn.Ping()
	if err != nil {
		logger.Critical("couldn't ping database: %s", err.Error())
	}

	conn.Exec("SET lc_time_names = 'es_ES';")
	_, err = conn.Exec("SET @@session.time_zone='+00:00';")
	if err != nil {
		logger.Critical("error configuring database: %s", err.Error())
	}
	logger.Information("database configured")
	db = conn
}

func GetDB() *sqlx.DB {
	if db == nil {
		start()
	}
	return db
}
