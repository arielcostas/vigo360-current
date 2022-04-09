/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
	"errors"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	req_post_id := mux.Vars(r)["postid"]

	var post, err = GetFullPost(req_post_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warning("[post] could not find post with that id")
			NotFoundHandler(w, r)
			return
		}
		logger.Error("[post] unexpected error fetching post from database: %s", err.Error())
	}

	// Fetch series
	var serie Serie
	if post.Serie.Valid {
		serie, err = GetSerieById(post.Serie.String)
		if err != nil {
			logger.Warning("[post] error fetching serie for post %s: %s", post.Id, err.Error())
			InternalServerErrorHandler(w, r)
			return
		}
	}

	recommendations, err := generateSuggestions(post.Id)
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(w, "post.html", struct {
		Post            FullPost
		Recommendations []PostRecommendation
		Meta            PageMeta
		Serie           Serie
	}{
		Serie:           serie,
		Post:            post,
		Recommendations: recommendations,
		Meta: PageMeta{
			Titulo:      post.Titulo,
			Descripcion: post.Resumen,
			Keywords:    post.Tags.String,
			Canonica:    FullCanonica("/post/" + post.Id),
			Miniatura:   FullCanonica("/static/thumb/" + post.Id + ".jpg"),
		},
	})
	if err != nil {
		logger.Error("[autores] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}
