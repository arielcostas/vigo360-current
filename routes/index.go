package routes

import (
	"fmt"
	"net/http"
)

type Resp struct {
	Now string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT NOW()")

	if err != nil {
		fmt.Println("error querying " + err.Error())
	}

	for rows.Next() {
		var r Resp
		rows.Scan(&r.Now)
		w.Write([]byte("row => " + r.Now))
	}
}
