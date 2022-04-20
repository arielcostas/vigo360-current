/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"encoding/xml"
	"net/http"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
)

type SitemapQuery struct {
	Uri                 string `xml:"loc"`
	Fecha_actualizacion string `xml:"lastmod,omitempty"`
	Changefreq          string `xml:"changefreq,omitempty"`
	Priority            string `xml:"priority,omitempty"`
}

type SitemapPage struct {
	XMLName xml.Name       `xml:"urlset"`
	Data    []SitemapQuery `xml:"url"`
}

func DatabaseAppError(err error) *appError {
	return &appError{err, "error reading data", "Error leyendo datos", 500}
}

func GenerateSitemap(w http.ResponseWriter, r *http.Request) *appError {
	var (
		pages = []SitemapQuery{}
		db    = database.GetDB()
		as    = model.NewAutorStore(db)
		tbs   = model.NewTrabajoStore(db)
		tas   = model.NewTagStore(db)
		ps    = model.NewPublicacionStore(db)
	)

	autores, err := as.Listar()
	if err != nil {
		return DatabaseAppError(err)
	}

	trabajos, err := tbs.Listar()
	if err != nil {
		return DatabaseAppError(err)
	}
	trabajos = trabajos.FiltrarPublicos()

	tags, err := tas.Listar()
	if err != nil {
		return DatabaseAppError(err)
	}

	publicaciones, err := ps.Listar()
	if err != nil {
		return DatabaseAppError(err)
	}
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
	for _, v := range nodbPageMeta {
		pages = append(pages, SitemapQuery{Uri: v.Canonica, Changefreq: "yearly", Priority: "0.3"})
	}

	output, err := xml.MarshalIndent(SitemapPage{Data: pages}, "", "\t")
	if err != nil {
		return &appError{err, "error marshalling xml", "Error produciendo página", 500}
	}
	w.Header().Add("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	w.Write(output)
	return nil
}
