package internal

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

type SitemapQuery struct {
	Uri                 string `xml:"loc"`
	Fecha_actualizacion string `xml:"lastmod,omitempty"`
	Changefreq          string `xml:"changefreq,omitempty"`
	Priority            string `xml:"priority,omitempty"`
}

type SitemapPage struct {
	XMLName xml.Name       `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Data    []SitemapQuery `xml:"url"`
}

func (s *Server) handlePublicSitemap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var pages = []SitemapQuery{}

		autores, err := s.store.autor.Listar()
		if err != nil {
			logger.Error("error recuperando autores: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		trabajos, err := s.store.trabajo.Listar()
		if err != nil {
			logger.Error("error recuperando trabajos: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		tags, err := s.store.tag.Listar()
		if err != nil {
			logger.Error("error recuperando tags: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		publicaciones, err := s.store.publicacion.Listar()
		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		trabajos = trabajos.FiltrarPublicos()
		publicaciones = publicaciones.FiltrarPublicas()

		pages = append(pages, SitemapQuery{Uri: "/", Changefreq: "daily", Priority: "0.8"})
		pages = append(pages, SitemapQuery{Uri: "/autores", Changefreq: "monthly", Priority: "0.5"})
		pages = append(pages, SitemapQuery{Uri: "/trabajos", Changefreq: "monthly", Priority: "0.5"})
		pages = append(pages, SitemapQuery{Uri: "/tags", Changefreq: "monthly", Priority: "0.5"})

		for _, autor := range autores {
			pages = append(pages, SitemapQuery{Uri: "/autores/" + autor.Id, Changefreq: "weekly", Priority: "0.3"})
		}
		for _, tag := range tags {
			pages = append(pages, SitemapQuery{Uri: "/tags/" + tag.Id, Changefreq: "weekly", Priority: "0.3"})
		}

		for _, trabajo := range trabajos {
			pages = append(pages, SitemapQuery{Uri: "/trabajos/" + trabajo.Id, Changefreq: "monthly", Priority: "0.3"})
		}
		for _, post := range publicaciones {
			pages = append(pages, SitemapQuery{Uri: "/post/" + post.Id, Changefreq: "monthly", Priority: "0.3"})
		}

		var domain = os.Getenv("DOMAIN")
		for i, p := range pages {
			p.Uri = domain + p.Uri
			pages[i] = p
		}

		// TODO: Mostrar páginas sin base de datos

		output, err := xml.MarshalIndent(SitemapPage{Data: pages}, "", "\t")
		if err != nil {
			logger.Error("error produciendo XML: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
		w.Header().Add("Content-Type", "application/xml")
		fmt.Fprintf(w, "%s%s", xml.Header, output)
	}
}
