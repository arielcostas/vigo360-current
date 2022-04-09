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

func PostPage(w http.ResponseWriter, r *http.Request) *appError {
	req_post_id := mux.Vars(r)["postid"]
	var post FullPost
	if np, err := GetFullPost(req_post_id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "post with that ID not found", Response: "No se ha encontrado una publicación en esta URL.", Status: 404}
		}
		return &appError{Error: err, Message: "database error fetching post", Response: "Error obteniendo ladatos.", Status: 500}
	} else {
		post = np
	}

	var serie Serie
	if post.Serie.Valid {
		if ns, err := GetSerieById(post.Serie.String); err != nil {
			return &appError{Error: err, Message: "database error fetching series for post", Response: "Error obteniendo datos.", Status: 500}
		} else {
			serie = ns
		}
	}

	var recommendations []PostRecommendation
	if nr, err := generateSuggestions(post.Id); err != nil {
		logger.Error("[%s] error fetching recommendations: %s", r.Context().Value("rid"), err.Error())
		recommendations = make([]PostRecommendation, 0)
	} else {
		recommendations = nr
	}

	var err = t.ExecuteTemplate(w, "post.html", struct {
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
		return &appError{Error: err, Message: "error rendering template", Response: "Hubo un error mostrando la página solicitada.", Status: 500}
	}

	return nil
}
