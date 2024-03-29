package internal

import (
	"math/rand"
	"net/http"
	"sort"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicListTags() http.HandlerFunc {
	type response struct {
		Tags []models.Tag
		Meta PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var tags, err = s.store.tag.Listar()
		for i, t := range tags {
			if t.Publicaciones < 1 {
				tags = append(tags[:i], tags[i+1:]...)
			}
		}

		if err != nil {
			logger.Error("error obteniendo tags: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
		}

		sort.Slice(tags, func(p, q int) bool {
			return tags[p].Nombre < tags[q].Nombre
		})

		/* Obtener última publicación para cada tag evitando duplicaciones */
		var publicacionesUsadas = make(map[string]bool, 0)
		for i, t := range tags {
			nt := t
			var publicacionesConTag []string
			pt, err := s.store.publicacion.ListarPorTag(t.Id)
			pt.FiltrarPublicas()
			pt.FiltrarRetiradas()
			for _, p := range pt {
				publicacionesConTag = append(publicacionesConTag, p.Id)
			}
			if err != nil {
				logger.Error("error obteniendo publicaciones con tag: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
				return
			}

			for _, pub := range publicacionesConTag {
				if _, ok := publicacionesUsadas[pub]; !ok {
					publicacionesUsadas[pub] = true
					nt.Ultima = pub
					break
				}
			}

			// Si se diera el caso de que todas están escogidas, poner una al azar
			if nt.Ultima == "" {
				cant := len(publicacionesConTag) - 1
				if cant <= 1 {
					cant = 1
				}
				aleatoria := rand.Intn(cant)
				nt.Ultima = publicacionesConTag[aleatoria]
			}

			tags[i] = nt
		}

		err = templates.Render(w, "tags.html", response{
			Tags: tags,
			Meta: PageMeta{
				Titulo:      "Tags",
				Descripcion: "Las diversas tags en las que se categorizan los artículos de Vigo360",
				Canonica:    fullCanonica("/tags"),
				BaseUrl:     baseUrl(),
			},
		})

		if err != nil {
			logger.Error("error generando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
