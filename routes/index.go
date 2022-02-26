package routes

import (
	"log"
	"net/http"
)

type IndexPost struct {
	Id                string
	Fecha_publicacion string
	Alt_portada       string
	Titulo            string
	Resumen           string
	Autor_id          string
	Autor_nombre      string `db:"nombre"`
}

type IndexQuery struct {
	Posts []IndexPost
}

type IndexParams struct {
	Posts []IndexPost
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	posts := []IndexPost{}
	err := db.Select(&posts, "SELECT publicaciones.id, DATE_FORMAT(publicaciones.fecha_publicacion, '%c %b.') as fecha_publicacion, publicaciones.alt_portada, publicaciones.titulo, publicaciones.resumen, autores.nombre FROM publicaciones LEFT JOIN autores on publicaciones.autor_id = autores.id;")
	if err != nil {
		log.Fatalf(err.Error())
	}

	t.ExecuteTemplate(w, "index.html", IndexParams{
		Posts: posts,
	})
}
