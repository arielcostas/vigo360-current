package public

import (
	"log"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	posts := []ResumenPost{}
	err := db.Select(&posts, "SELECT publicaciones.id, DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, publicaciones.alt_portada, publicaciones.titulo, publicaciones.resumen, autores.nombre FROM publicaciones LEFT JOIN autores on publicaciones.autor_id = autores.id WHERE publicaciones.fecha_publicacion < NOW() ORDER BY publicaciones.fecha_publicacion DESC;")
	if err != nil {
		log.Fatalf(err.Error())
	}

	t.ExecuteTemplate(w, "index.html", struct {
		Posts []ResumenPost
		Meta  common.PageMeta
	}{
		Posts: posts,
		Meta: common.PageMeta{
			Titulo:      "Inicio",
			Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
			Canonica:    FullCanonica("/"),
		},
	})
}
