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

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/model"
)

func PostPage(w http.ResponseWriter, r *http.Request) *appError {
	req_post_id := mux.Vars(r)["postid"]
	var (
		db = database.GetDB()
		ps = model.NewPublicacionStore(db)
		ss = model.NewSerieStore(db)
		e2 error
	)
	var post model.Publicacion
	if np, err := ps.ObtenerPorId(req_post_id, true); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "post with that ID not found", Response: "No se ha encontrado una publicación en esta URL.", Status: 404}
		}
		return &appError{Error: err, Message: "database error fetching post", Response: "Error obteniendo datos.", Status: 500}
	} else {
		post = np
	}

	if post.Serie.Id != "" {
		post.Serie, e2 = ss.Obtener(post.Serie.Id)
		if e2 != nil {
			return &appError{Error: e2, Message: "database error fetching series for post", Response: "Error obteniendo datos.", Status: 500}
		}
	}

	var recommendations []Sugerencia
	if nr, err := generateSuggestions(post.Id); err != nil {
		logger.Error("[%s] error fetching recommendations: %s", r.Context().Value("rid"), err.Error())
		recommendations = make([]Sugerencia, 0)
	} else {
		recommendations = nr
	}

	var keywords = ""
	for _, t := range post.Tags {
		keywords += t.Nombre + ","
	}

	var output bytes.Buffer
	var err = t.ExecuteTemplate(&output, "post.html", struct {
		Post            model.Publicacion
		Recommendations []Sugerencia
		Meta            PageMeta
	}{
		Post:            post,
		Recommendations: recommendations,
		Meta: PageMeta{
			Titulo:      post.Titulo,
			Descripcion: post.Resumen,
			Keywords:    keywords,
			Canonica:    FullCanonica("/post/" + post.Id),
			Miniatura:   FullCanonica("/static/thumb/" + post.Id + ".jpg"),
		},
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Hubo un error mostrando la página solicitada.", Status: 500}
	}

	w.Write(output.Bytes())
	return nil
}
