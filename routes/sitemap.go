package routes

import (
	"log"
	"net/http"
)

type SitemapQuery struct {
	Uri                 string
	Fecha_actualizacion string
	Changefreq          string
	Priority            string
}

type SitemapPage struct {
	Urls []SitemapQuery
}

func GenerateSitemap(w http.ResponseWriter, r *http.Request) {
	pages := []SitemapQuery{}
	query := `SELECT * FROM sitemap;`

	err := db.Select(&pages, query)
	if err != nil {
		log.Fatalf(err.Error())
	}

	pages = append(pages, SitemapQuery{Uri: "/licencias", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/contacto", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/siguenos", Changefreq: "yearly", Priority: "0.3"})

	tt.ExecuteTemplate(w, "sitemap.xml", SitemapPage{
		Urls: pages,
	})
}
