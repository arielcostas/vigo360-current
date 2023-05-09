package internal

import (
	"net/http"
	"strings"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicBusqueda() http.HandlerFunc {
	type resultado struct {
		Id      string
		Titulo  string
		Resumen string
		Uri     string
	}

	type response struct {
		Resultados []resultado
		Termino    string
		Meta       PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var resultados = make([]resultado, 0)

		var termino = r.URL.Query().Get("termino")
		termino = strings.TrimSpace(termino)

		if termino == "" {
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
		}

		autores, err := s.store.autor.Buscar(termino)
		if err != nil {
			logger.Error("error obteniendo autores: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		for _, autor := range autores {
			resultados = append(resultados, resultado{
				Id:      autor.Id,
				Titulo:  autor.Nombre,
				Resumen: autor.Biografia,
				Uri:     "/autores/" + autor.Id,
			})
		}

		publicaciones, err := s.store.publicacion.Buscar(termino)
		publicaciones = publicaciones.FiltrarPublicas()
		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		for _, pub := range publicaciones {
			resultados = append(resultados, resultado{
				Id:      pub.Id,
				Titulo:  pub.Titulo,
				Resumen: pub.Resumen,
				Uri:     "/post/" + pub.Id,
			})
		}

		err = templates.Render(w, "search.html", response{
			Resultados: resultados,
			Termino:    termino,
			Meta: PageMeta{
				Titulo:   "Resultados para " + termino,
				Canonica: fullCanonica("/buscar?termino=" + termino),
			},
		})
		if err != nil {
			logger.Error("error generando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
			return
		}
	}
}
