/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"net/http"

	"vigo360.es/new/internal/database"
)

type SitemapQuery struct {
	Uri                 string `xml:"loc"`
	Fecha_actualizacion string `xml:"lastmod"`
	Changefreq          string `xml:"changefreq"`
	Priority            string `xml:"priority"`
}

type SitemapPage struct {
	XMLName xml.Name       `xml:"urlset"`
	Data    []SitemapQuery `xml:"url"`
}

func GenerateSitemap(w http.ResponseWriter, r *http.Request) *appError {
	pages := []SitemapQuery{}
	query := `SELECT * FROM sitemap;`

	db := database.GetDB()
	err := db.Select(&pages, query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{err, "error fetching rows", "Error obteniendo datos", 500}
	}

	pages = append(pages, SitemapQuery{Uri: "/licencias", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/contacto", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/siguenos", Changefreq: "yearly", Priority: "0.3"})

	output, err := xml.MarshalIndent(SitemapPage{Data: pages}, "", "\t")
	if err != nil {
		return &appError{err, "error marshalling xml", "Error produciendo página", 500}
	}
	w.Header().Add("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	w.Write(output)
	return nil
}
