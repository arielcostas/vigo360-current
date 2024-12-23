package internal

import (
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicNodbPage() http.HandlerFunc {
	var nodbPageMeta = map[string]PageMeta{
		"policy": {
			Titulo:      "Políticas de Vigo360",
			Descripcion: "Información sobre políticas relativa a Vigo360: uso de contenidos, privacidad, rectificación...",
			Canonica:    fullCanonica("/policy"),
			BaseUrl:     baseUrl(),
		},
		"contacto": {
			Titulo:      "Contacto",
			Descripcion: "Si necesitases contactar con Vigo360, aquí encontrarás cómo hacerlo.",
			Canonica:    fullCanonica("/contacto"),
			BaseUrl:     baseUrl(),
		},
	}

	type response struct {
		Meta PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var (
			page = r.URL.Path[1:]
			meta PageMeta
		)

		if pm, ok := nodbPageMeta[page]; ok {
			meta = pm
		} else {
			logger.Error("nodb page not found")
			s.handleError(r, w, 404, messages.ErrorPaginaNoEncontrada)
			return
		}

		err := templates.Render(w, page+".html", response{
			Meta: meta,
		})
		if err != nil {
			logger.Error("error mostrando la página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}

}
