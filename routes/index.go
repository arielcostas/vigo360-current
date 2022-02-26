package routes

import (
	"net/http"
)

type Resp struct {
	Now string
}

type IndexParams struct {
	Row string
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	row := db.QueryRow("SELECT curtime()")
	var resp Resp = Resp{}
	row.Scan(&resp.Now)

	t.ExecuteTemplate(w, "index.html", &IndexParams{
		Row: resp.Now,
	})
}
