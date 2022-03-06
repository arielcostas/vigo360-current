package public

import (
	"log"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"github.com/gorilla/mux"
)

func TagsIdPage(w http.ResponseWriter, r *http.Request) {
	req_tagid := mux.Vars(r)["tagid"]

	tag := Tag{}
	err := db.QueryRowx("SELECT nombre FROM tags WHERE id=?;", req_tagid).StructScan(&tag)
	if err != nil {
		log.Fatalf(err.Error())
	}

	posts := []ResumenPost{}
	db.Select(&posts, `SELECT publicaciones.id, DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, publicaciones.alt_portada, publicaciones.titulo,
	autores.nombre FROM publicaciones_tags
	LEFT JOIN publicaciones ON publicaciones_tags.publicacion_id = publicaciones.id
    LEFT JOIN autores ON publicaciones.autor_id = autores.id
    WHERE tag_id = ? ORDER BY publicaciones.fecha_publicacion DESC;`, req_tagid)

	t.ExecuteTemplate(w, "tags-id.html", struct {
		Tag   Tag
		Posts []ResumenPost
		Meta  common.PageMeta
	}{
		Tag:   tag,
		Posts: posts,
		Meta: common.PageMeta{
			Titulo:      tag.Nombre,
			Descripcion: "Publicaciones en Vigo360 sobre " + tag.Nombre,
			Canonica:    FullCanonica("/tags/" + req_tagid),
		},
	})
}
