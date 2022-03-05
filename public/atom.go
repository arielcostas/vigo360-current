package public

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type PostsAtomPost struct {
	Id                  string
	Fecha_publicacion   string
	Publicacion_3339    string
	Fecha_actualizacion string
	Actualizacion_3339  string

	Titulo       string
	Resumen      string
	Autor_id     string
	Autor_nombre string
	Autor_email  string
	Tag_id       string
	Raw_tags     string
	Tags         []string
}

type PostsAtomTag struct {
	Id     string
	Nombre string
}

type PostsAtomParams struct {
	BaseURL string
	Now     string
	Posts   []PostsAtomPost
}

func PostsAtomFeed(w http.ResponseWriter, r *http.Request) {
	tags := []PostsAtomTag{}
	tagMap := map[string]string{}
	err := db.Select(&tags, `SELECT * FROM tags`)

	for _, tag := range tags {
		tagMap[tag.Id] = tag.Nombre
	}

	if err != nil {
		log.Fatalf(err.Error())
	}

	posts := []PostsAtomPost{}
	err = db.Select(&posts, `SELECT publicaciones.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, tag_id, GROUP_CONCAT(tag_id) as raw_tags FROM publicaciones LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id LEFT JOIN autores ON publicaciones.autor_id = autores.id GROUP BY id;`)

	if err != nil {
		log.Fatalf(err.Error())
	}

	for i := 0; i < len(posts); i++ {
		p := posts[i]
		t, _ := time.Parse("2006-01-02 15:04:05", p.Fecha_publicacion)
		p.Publicacion_3339 = t.Format(time.RFC3339)

		t, _ = time.Parse("2006-01-02 15:04:05", p.Fecha_actualizacion)
		p.Actualizacion_3339 = t.Format(time.RFC3339)

		p.Id = url.PathEscape(p.Id)
		p.Tags = []string{}

		for _, tag := range strings.Split(p.Raw_tags, ",") {
			p.Tags = append(p.Tags, tagMap[tag])
		}

		posts[i] = p
	}

	w.Header().Add("Content-Type", "application/atom+xml")
	err = tt.ExecuteTemplate(w, "atom.xml", PostsAtomParams{
		BaseURL: os.Getenv("DOMAIN"),
		Now:     time.Now().Format("2006-01-02T15:04:05-07:00"),
		Posts:   posts,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}
}
