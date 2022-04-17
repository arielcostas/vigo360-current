/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"

	"vigo360.es/new/internal/logger"
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
	if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[sitemap]: unable to fetch rows: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	pages = append(pages, SitemapQuery{Uri: "/licencias", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/contacto", Changefreq: "yearly", Priority: "0.3"})
	pages = append(pages, SitemapQuery{Uri: "/siguenos", Changefreq: "yearly", Priority: "0.3"})

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "sitemap.xml", SitemapPage{
		Urls: pages,
	})
	if err != nil {
		logger.Error("[sitemap] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
	w.Write(output.Bytes())
}
