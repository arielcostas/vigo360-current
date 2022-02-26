package routes

import (
	"log"
	"net/http"
)

type IndexParams struct {
	Row string
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	row := db.QueryRow("SELECT curtime() as Now")
	dbresp := struct{ Now string }{}
	err := row.Scan(&dbresp.Now)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println(dbresp.Now)

	t.ExecuteTemplate(w, "index.html", &IndexParams{
		Row: dbresp.Now,
	})
}
