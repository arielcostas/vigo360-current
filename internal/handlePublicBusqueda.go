package internal

import (
	"net/http"
	"os"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicBusqueda() http.HandlerFunc {
	type resultado struct {
		Id                string
		Titulo            string
		Resumen           string
		Autor_nombre      string
		Alt_portada       string
		Uri               string
		Fecha_publicacion string
	}

	type response struct {
		Resultados []resultado
		Termino    string
		Meta       PageMeta
	}

	var algoliaApiKey string = os.Getenv("ALGOLIA_API_KEY")
	var algoliaAppId string = os.Getenv("ALGOLIA_APP_ID")
	var algoliaIndexName string = os.Getenv("ALGOLIA_INDEX_NAME")

	var searchOptions = []interface{}{
		opt.Filters("fecha_publicacion != null"),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var resultados = make([]resultado, 0)

		var termino = r.URL.Query().Get("termino")
		termino = strings.TrimSpace(termino)

		if termino == "" {
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
		}

		client := search.NewClient(algoliaAppId, algoliaApiKey)
		index := client.InitIndex(algoliaIndexName)

		result, err := index.Search(termino, searchOptions)
		if err != nil {
			log.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		for _, hit := range result.Hits {
			resultados = append(resultados, resultado{
				Id:                hit["id"].(string),
				Titulo:            hit["titulo"].(string),
				Autor_nombre:      hit["autor_nombre"].(string),
				Alt_portada:       hit["alt_portada"].(string),
				Resumen:           hit["resumen"].(string),
				Uri:               "/post/" + hit["id"].(string),
				Fecha_publicacion: hit["fecha_publicacion"].(string),
			})
		}

		//publicaciones, err := s.store.publicacion.Buscar(termino)
		//publicaciones = publicaciones.FiltrarPublicas().FiltrarRetiradas()
		//if err != nil {
		//	log.Error("error recuperando publicaciones: %s", err.Error())
		//	s.handleError(r, w, 500, messages.ErrorDatos)
		//	return
		//}

		//for _, pub := range publicaciones {
		//	resultados = append(resultados, resultado{
		//		Id:                pub.Id,
		//		Titulo:            pub.Titulo,
		//		Autor_nombre:      pub.Autor.Nombre,
		//		Alt_portada:       pub.Alt_portada,
		//		Resumen:           pub.Resumen,
		//		Uri:               "/post/" + pub.Id,
		//		Fecha_publicacion: pub.Fecha_publicacion,
		//	})
		//}

		err = templates.Render(w, "search.html", response{
			Resultados: resultados,
			Termino:    termino,
			Meta: PageMeta{
				Titulo:   "Resultados para " + termino,
				Canonica: fullCanonica("/buscar?termino=" + termino),
				BaseUrl:  baseUrl(),
			},
		})
		if err != nil {
			log.Error("error generando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
			return
		}
	}
}
