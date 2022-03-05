package public

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

type IndexParams struct {
	Posts []IndexPost
	Meta  PageMeta
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	posts := []IndexPost{}
	err := db.Select(&posts, "SELECT publicaciones.id, DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, publicaciones.alt_portada, publicaciones.titulo, publicaciones.resumen, autores.nombre FROM publicaciones LEFT JOIN autores on publicaciones.autor_id = autores.id WHERE publicaciones.fecha_publicacion < NOW() ORDER BY publicaciones.fecha_publicacion DESC;")
	if err != nil {
		log.Fatalf(err.Error())
	}

	t.ExecuteTemplate(w, "index.html", IndexParams{
		Posts: posts,
		Meta: PageMeta{
			Titulo:      "Inicio",
			Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
			Canonica:    FullCanonica("/"),
		},
	})
}
