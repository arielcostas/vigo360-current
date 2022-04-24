package public

import (
	"bytes"
	"net/http"
	"strings"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
)

func realizarBusqueda(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		ps = model.NewPublicacionStore(db)
	)

	var termino = r.URL.Query().Get("termino")
	termino = strings.TrimSpace(termino)
	// TODO Gestionar-impedir términos vacíos
	if termino == "" {
		w.Header().Add("Location", "/")
		w.WriteHeader(302)
		return nil
	}

	resultados, err := ps.Buscar(termino)
	resultados = resultados.FiltrarPublicas()
	if err != nil {
		return &appError{err, "error searching", "Hubo un error realizando la búsqueda", 500}
	}

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "search.html", struct {
		Publicaciones model.Publicaciones
		Termino       string
		Meta          PageMeta
	}{
		Publicaciones: resultados,
		Termino:       termino,
		Meta: PageMeta{
			Titulo:   "Resultados para " + termino,
			Canonica: FullCanonica("/buscar?termino=" + termino),
		},
	})
	if err != nil {
		return &appError{err, "error rendering template", "Hubo un error mostrando la página solicitada", 500}
	}
	w.Write(output.Bytes())
	return nil
}
