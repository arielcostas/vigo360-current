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
	"github.com/go-playground/validator/v10"
)

type ResumenPost struct {
	Id                string
	Titulo            string
	Fecha_publicacion sql.NullString
	CantTags          int
	Publicado         bool
	Autor_id          string
	Autor_nombre      string
}

//go:embed extra/default.jpg
var defaultImageJPG []byte

//go:embed extra/default.webp
var defaultImageWebp []byte

func listPosts(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	_, err = getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	posts := []ResumenPost{}

	err = db.Select(&posts, `SELECT publicaciones.id, titulo, (fecha_publicacion < NOW() && fecha_publicacion IS NOT NULL) as publicado, DATE_FORMAT(fecha_publicacion,'%d-%m-%Y') as fecha_publicacion, autor_id, autores.nombre as autor_nombre, count(tag_id) as canttags FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id GROUP BY publicaciones.id ORDER BY publicado ASC, publicaciones.fecha_publicacion DESC;`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "posts")
	}

	err = t.ExecuteTemplate(w, "post.html", struct {
		Posts []ResumenPost
	}{
		Posts: posts,
	})

	if err != nil {
		return newTemplateRenderingAppError(err)
	}

	return nil
}

// Data to be input via a form to create a new post
func createPost(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	sess, err := getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	var art_autor = sess.Autor_id

	if err := r.ParseForm(); err != nil {
		return &appError{Error: err, Message: "error parsing form",
			Response: "Hubo un error leyendo el formulario", Status: 500}
	}

	fi := struct {
		Id     string `validate:"required,min=3,max=40,lowercase"`
		Titulo string `validate:"required,min=3,max=80"`
	}{}
	fi.Id = r.FormValue("art-id")
	fi.Titulo = r.FormValue("art-titulo")

	if err := validator.New().Struct(fi); err != nil {
		// TODO: Mostrar de nuevo formulario con los errores
		return &appError{Error: err, Message: "error de validacion",
			Response: "Ha habido un error validando el formulario. Revise los campos e inténtelo de nuevo.", Status: 400}
	}

	var tx *sql.Tx
	if nt, err := db.Begin(); err != nil {
		return &appError{Error: err, Message: "error beginning transaction",
			Response: "Error creando publicación en la de datos", Status: 500}
	} else {
		tx = nt
	}

	q := "INSERT INTO publicaciones(id, titulo, alt_portada, resumen, contenido, autor_id) VALUES (?, ?, 'CAMBIAME','', '', ?);"
	if _, err := tx.Exec(q, fi.Id, fi.Titulo, art_autor); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error creating article in database",
			Response: "Error creando publicación en la de datos", Status: 500}
	}

	photopath := os.Getenv("UPLOAD_PATH")
	if err := os.WriteFile(photopath+"/images/"+fi.Id+".webp", defaultImageWebp, 0o644); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error saving webp to disk",
			Response: "Hubo un error guardando el artículo", Status: 500}
	}

	if err := os.WriteFile(photopath+"/thumb/"+fi.Id+".jpg", defaultImageJPG, 0o644); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error saving jpg to disk",
			Response: "Hubo un error guardando el artículo", Status: 500}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error performing transaction",
			Response: "Error creando publicación en la de datos", Status: 500}
	}

	w.Header().Add("Location", "/admin/post/"+fi.Id)
	w.WriteHeader(303)
	return nil
}
