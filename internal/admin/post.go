/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"bytes"
	"database/sql"
	_ "embed"
	"errors"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
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
	sess, err := getSession(sc.Value)
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

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "post.html", struct {
		Posts   []ResumenPost
		Session Session
	}{
		Posts:   posts,
		Session: sess,
	})

	if err != nil {
		return newTemplateRenderingAppError(err)
	}

	w.Write(output.Bytes())
	return nil
}

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

	// TODO: Refactor this
	fi := struct {
		Titulo string `validate:"required,min=3,max=80"`
	}{
		Titulo: r.FormValue("art-titulo"),
	}

	var id = r.FormValue("art-id")

	if !ValidatePostId(id) {
		return &appError{Error: ErrInvalidFormInput, Message: "error de validacion",
			Response: "Ha habido un error validando el formulario. Revise los campos e inténtelo de nuevo.", Status: 400}
	}

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
	if _, err := tx.Exec(q, id, fi.Titulo, art_autor); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error creating article in database",
			Response: "Error creando publicación en la de datos", Status: 500}
	}

	photopath := os.Getenv("UPLOAD_PATH")
	if err := os.WriteFile(photopath+"/images/"+id+".webp", defaultImageWebp, 0o644); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error saving webp to disk",
			Response: "Hubo un error guardando el artículo", Status: 500}
	}

	if err := os.WriteFile(photopath+"/thumb/"+id+".jpg", defaultImageJPG, 0o644); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error saving jpg to disk",
			Response: "Hubo un error guardando el artículo", Status: 500}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error performing transaction",
			Response: "Error creando publicación en la de datos", Status: 500}
	}

	w.Header().Add("Location", "/admin/post/"+id)
	w.WriteHeader(303)
	return nil
}

func deletePost(w http.ResponseWriter, r *http.Request) *appError {
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
	if !sess.Permisos["publicaciones_delete"] {
		return &appError{Error: ErrUnablePermissions, Message: "user doesn't have permission to delete posts", Response: "No tienes permiso para realizar esta acción", Status: 403}
	}

	var postid = mux.Vars(r)["postid"]
	tx, err := db.Begin()
	if err != nil {
		return &appError{Error: err, Message: "error beginning transaction", Response: "Error eliminando la publicación", Status: 500}
	}
	_, err = tx.Exec("DELETE FROM publicaciones_tags WHERE publicacion_id=?", postid)
	if err != nil {
		return &appError{Error: err, Message: "error deleting from publicaciones_tags", Response: "Error eliminando la publicación", Status: 500}
	}
	_, err = tx.Exec("DELETE FROM publicaciones WHERE id=?", postid)
	if err != nil {
		return &appError{Error: err, Message: "error deleting from publicaciones", Response: "Error eliminando la publicación", Status: 500}
	}
	err = tx.Commit()
	if err != nil {
		return &appError{Error: err, Message: "error committing transaction", Response: "Error eliminando la publicación", Status: 500}
	}

	w.Header().Add("Location", "/admin/post")
	defer w.WriteHeader(307)
	return nil
}
