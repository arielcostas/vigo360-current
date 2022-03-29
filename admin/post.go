/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"database/sql"
	_ "embed"
	"errors"
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

type ResumenPost struct {
	Id           string
	Titulo       string
	Publicado    bool
	Autor_id     string
	Autor_nombre string
}

//go:embed extra/default.jpg
var defaultImageJPG []byte

//go:embed extra/default.webp
var defaultImageWebp []byte

func PostListPage(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	posts := []ResumenPost{}

	err := db.Select(&posts, `SELECT publicaciones.id, titulo, (fecha_publicacion < NOW() && fecha_publicacion IS NOT NULL) as publicado, autor_id, autores.nombre as autor_nombre FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id ORDER BY publicado ASC, fecha_publicacion DESC;`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("error inesperado leyendo publicaciones de la base de datos: %s", err.Error())
	}

	t.ExecuteTemplate(w, "post.html", struct {
		Posts []ResumenPost
	}{
		Posts: posts,
	})
}

// HTTP Handler for creating posts, accessible by authenticated authors via `POST /admin/post`. It requires passing art-id and art-titulo in the body, as part of the form submission from the `GET` page with the same URI.
func CreatePostAction(w http.ResponseWriter, r *http.Request) {
	sesion := verifyLogin(w, r)
	err := r.ParseForm()
	if err != nil {
		logger.Error("error parsing create-post form: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	art_id := r.FormValue("art-id")
	art_titulo := r.FormValue("art-titulo")
	art_autor := sesion.Id

	// Article id must be below 40 characters long, with only lowercase spanish letters, numbers and dashes
	if !validarId(art_id) {
		w.WriteHeader(400)
		// TODO: proper error page
		w.Write([]byte("El id del artículo debe contener entre 3 y 40 letras minúsculas del alfabeto español, números, guiones o guiones bajos."))
		return
	}

	if !validarTitulo(art_titulo) {
		w.WriteHeader(400)
		w.Write([]byte("El título tiene que contener entre 3 y 80 caracteres."))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(500)
		logger.Error("[post] error beginning insert operation: %s", err.Error())
		w.Write([]byte("Error creando el artículo"))
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	_, err = tx.Exec(`INSERT INTO publicaciones(id, titulo, alt_portada, resumen, contenido, autor_id) VALUES (?, ?, "CAMBIAME","", "", ?);`, art_id, art_titulo, art_autor)

	if err != nil {
		// TODO: proper error page
		w.WriteHeader(500)
		logger.Error("[post] error creating article in database: %s", err.Error())
		w.Write([]byte("Error creando el artículo"))
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	err = tx.Commit()

	if err != nil {
		// TODO: proper error page
		w.WriteHeader(500)
		logger.Error("[post] error commiting article in database: %s", err.Error())
		w.Write([]byte("Error creando el artículo"))
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	// Every article needs its default photo
	photopath := os.Getenv("UPLOAD_PATH")
	err = os.WriteFile(photopath+"/images/"+art_id+".webp", defaultImageWebp, 0o644)
	if err != nil {
		// TODO: proper error page
		w.WriteHeader(500)
		logger.Error("[post] error saving article webp: %s", err.Error())
		w.Write([]byte("Error creando foto WEBP predeterminada"))
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}
	err = os.WriteFile(photopath+"/thumb/"+art_id+".jpg", defaultImageJPG, 0o644)
	if err != nil {
		// TODO: proper error page
		w.WriteHeader(500)
		logger.Error("[post] error saving article jpg: %s", err.Error())
		w.Write([]byte("Error creando foto JPG predeterminada"))
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	w.Header().Add("Location", "/admin/post/"+art_id)
	w.WriteHeader(303)
}
