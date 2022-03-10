package public

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"github.com/gorilla/mux"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	req_post_id := mux.Vars(r)["postid"]
	query := `SELECT publicaciones.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_publicacion,
	DATE_FORMAT(publicaciones.fecha_actualizacion, '%e %b.') as fecha_actualizacion,
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
FROM publicaciones 
LEFT JOIN autores on publicaciones.autor_id = autores.id 
WHERE publicaciones.id = ? AND publicaciones.fecha_publicacion IS NOT NULL AND publicaciones.fecha_publicacion < NOW() ORDER BY publicaciones.fecha_publicacion DESC;`

	post := FullPost{}
	err := db.QueryRowx(query, req_post_id).StructScan(&post)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Result is in markdown, convert to HTML
	var buf bytes.Buffer
	common.Parser.Convert([]byte(post.ContenidoRaw), &buf)

	post.Contenido = template.HTML(buf.Bytes())

	t.ExecuteTemplate(w, "post.html", struct {
		Post FullPost
		Meta common.PageMeta
	}{
		Post: post,
		Meta: common.PageMeta{
			Titulo:      post.Titulo,
			Descripcion: post.Resumen,
			Canonica:    FullCanonica("/post/" + post.Id),
			Miniatura:   FullCanonica("/static/thumb/" + post.Id + ".jpg"),
		},
	})
}
