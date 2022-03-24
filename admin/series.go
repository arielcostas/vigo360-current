package admin

import (
	"net/http"
	"strings"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func ListSeries(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	series := []Serie{}
	err := db.Select(&series, `SELECT series.*, COUNT(publicaciones.id) as articulos FROM series LEFT JOIN publicaciones ON series.id = publicaciones.serie_id GROUP BY series.id;`)
	if err != nil {
		logger.Error("[series]: error fetching series from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	err = t.ExecuteTemplate(w, "series.html", struct {
		Series []Serie
	}{
		Series: series,
	})
}

func CreateSeries(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	err := r.ParseForm()

	if err != nil {
		logger.Error("[create-series] error parsing form: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	titulo := r.FormValue("titulo")
	id := strings.ToLower(strings.TrimSpace(titulo))
	id = strings.ReplaceAll(id, " ", "-")

	_, err = db.Exec(`INSERT INTO series VALUES (?, ?)`, id, titulo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error guardando nueva serie en la base de datos"))
		logger.Error("[create-series] error saving new series to database: %s", err.Error())
		return
	}

	w.Header().Add("Location", "/admin/series")
	w.WriteHeader(303)
}
