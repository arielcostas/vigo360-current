package public

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AtomPost struct {
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
	Tag_id       sql.NullString
	Raw_tags     sql.NullString
	Tags         []string
}

type FeedParams struct {
	BaseURL string
	Now     string
	Posts   []AtomPost
}

func PostsAtomFeed(w http.ResponseWriter, r *http.Request) {
	tags := []Tag{}
	tagMap := map[string]string{}
	err := db.Select(&tags, `SELECT * FROM tags`)

	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, tag := range tags {
		tagMap[tag.Id] = tag.Nombre
	}

	posts := []AtomPost{}
	err = db.Select(&posts, `SELECT publicaciones.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, tag_id, GROUP_CONCAT(tag_id) as raw_tags FROM publicaciones LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id LEFT JOIN autores ON publicaciones.autor_id = autores.id GROUP BY id;`)

	if err != nil {
		log.Fatalf(err.Error())
	}

	for i := 0; i < len(posts); i++ {
		p := posts[i]
		p.Tags = []string{}

		for _, tag := range strings.Split(p.Raw_tags.String, ",") {
			p.Tags = append(p.Tags, tagMap[tag])
		}

		posts[i] = p
	}

	writeFeed(w, r, "atom.xml", posts)
}

func TrabajosAtomFeed(w http.ResponseWriter, r *http.Request) {
	trabajos := []AtomPost{}
	err := db.Select(&trabajos, `SELECT trabajos.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM trabajos LEFT JOIN autores ON trabajos.autor_id = autores.id;`)

	if err != nil {
		log.Fatalf(err.Error())
	}

	writeFeed(w, r, "trabajos-atom.xml", trabajos)
}

func writeFeed(w http.ResponseWriter, r *http.Request, feedName string, items []AtomPost) {
	var lastUpdate time.Time

	for i := 0; i < len(items); i++ {
		p := items[i]
		t, _ := time.Parse("2006-01-02 15:04:05", p.Fecha_publicacion)
		p.Publicacion_3339 = t.Format(time.RFC3339)

		t, _ = time.Parse("2006-01-02 15:04:05", p.Fecha_actualizacion)
		p.Actualizacion_3339 = t.Format(time.RFC3339)

		if lastUpdate.Before(t) {
			lastUpdate = t
		}

		p.Id = url.PathEscape(p.Id)
	}

	w.Header().Add("Content-Type", "application/atom+xml")
	err := tt.ExecuteTemplate(w, feedName, &FeedParams{
		BaseURL: os.Getenv("DOMAIN"),
		Now:     lastUpdate.Format("2006-01-02T15:04:05-07:00"),
		Posts:   items,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}
}
