package admin

import (
	"database/sql"
	_ "embed"
	"errors"
	"net/http"
	"os"
	"regexp"

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
	if !regexp.MustCompile(`^[a-z|ñ|\-|\_|0-9]{3,40}$`).MatchString(art_id) {
		// TODO add a proper error page
		w.Write([]byte("El id del artículo debe contener entre 3 y 40 letras minúsculas del alfabeto español, números, guiones o guiones bajos."))
		return
	}

	if len(art_titulo) < 3 || len(art_titulo) > 80 {
		w.Write([]byte("El título tiene que contener entre 3 y 80 caracteres."))
		return
	}

	_, err = db.Exec(`INSERT INTO publicaciones(id, titulo, alt_portada, resumen, contenido, autor_id) VALUES (?, ?, "CAMBIAME","", "", ?);`, art_id, art_titulo, art_autor)

	if err != nil {
		// TODO add proper error page
		w.Write([]byte("Error creando el artículo"))
		logger.Error("error creating article in database: %s", err.Error())
		return
	}

	// Every article needs its default photo
	photopath := os.Getenv("UPLOAD_PATH")
	err = os.WriteFile(photopath+"/images/"+art_id+".webp", defaultImageWebp, 0o644)
	if err != nil {
		// TODO proper error page
		logger.Error("error creating article webp: %s", err.Error())
		w.Write([]byte("Error creating default WEBP photo"))
		return
	}
	err = os.WriteFile(photopath+"/thumb/"+art_id+".jpg", defaultImageJPG, 0o644)
	if err != nil {
		// TODO proper error page
		logger.Error("error creating article jpg: %s", err.Error())
		w.Write([]byte("Error creating default JPG photo"))
		return
	}

	w.Header().Add("Location", "/admin/post/"+art_id)
	w.WriteHeader(303)
}
