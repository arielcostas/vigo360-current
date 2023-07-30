package internal

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

type indexParams struct {
	CurrentPage int
	PageCount   int
	HasNextPage bool
	HasPrevPage bool
	IsFirstPage bool
	IsLastPage  bool
	Posts       models.Publicaciones
	Meta        PageMeta
}

func (s *Server) handlePublicIndex() http.HandlerFunc {
	var meta = PageMeta{
		Titulo:      "Inicio",
		Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
		Canonica:    fullCanonica("/"),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		posts, err := s.store.publicacion.Listar()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Error("error recuperando datos: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		posts = posts.FiltrarPublicas()

		/* Paginación */
		var pagina = 1

		var queryPage = r.URL.Query().Get("page")
		if queryPage != "" {
			o, err := strconv.Atoi(queryPage)
			if err != nil {
				log.Error("no se pudo convertir '%s' a un número de página", queryPage)
				s.handleError(r, w, 404, messages.ErrorNoResultados)
				return
			}
			pagina = o
		}

		var inicio = pagina*9 - 9
		var limite = getMinimo(inicio+9, len(posts))

		if inicio >= len(posts) || inicio < 0 {
			log.Error("con %d publicaciones no existe la página %s", len(posts), pagina)
			s.handleError(r, w, 404, messages.ErrorNoResultados)
			return
		}

		var restantes = len(posts) - 9 // Los artículos que aún no se han metido en una página
		var cantidadPaginas = 1
		for restantes > 0 {
			cantidadPaginas++
			restantes -= 9
		}

		err = templates.Render(w, "index.html", indexParams{
			CurrentPage: pagina,
			PageCount:   cantidadPaginas,
			Posts:       posts[inicio:limite],
			Meta:        meta,
		})

		if err != nil {
			log.Error("error renderizando la página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
