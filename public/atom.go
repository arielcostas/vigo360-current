package public

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

type AtomPost struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string

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
	BaseURL    string
	Id         string
	Nombre     string
	LastUpdate string
	Posts      []AtomPost
}

func PostsAtomFeed(w http.ResponseWriter, r *http.Request) {
	tags := []Tag{}
	tagMap := map[string]string{}
	err := db.Select(&tags, `SELECT * FROM tags`)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting tags: " + err.Error())
	}

	for _, tag := range tags {
		tagMap[tag.Id] = tag.Nombre
	}

	posts := []AtomPost{}
	err = db.Select(&posts, `SELECT publicaciones.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, tag_id, GROUP_CONCAT(tag_id) as raw_tags FROM publicaciones LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id LEFT JOIN autores ON publicaciones.autor_id = autores.id WHERE fecha_publicacion < NOW() GROUP BY id;`)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting posts: " + err.Error())
	}

	for i := 0; i < len(posts); i++ {
		p := posts[i]
		p.Tags = []string{}

		for _, tag := range strings.Split(p.Raw_tags.String, ",") {
			p.Tags = append(p.Tags, tagMap[tag])
		}

		posts[i] = p
	}

	writeFeed(w, r, "atom.xml", posts, "Publicaciones", "")
}

func TrabajosAtomFeed(w http.ResponseWriter, r *http.Request) {
	trabajos := []AtomPost{}
	err := db.Select(&trabajos, `SELECT trabajos.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM trabajos LEFT JOIN autores ON trabajos.autor_id = autores.id WHERE fecha_publicacion < NOW();`)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting trabajos: " + err.Error())
	}

	writeFeed(w, r, "trabajos-atom.xml", trabajos, "Trabajos", "")
}

func TagsAtomFeed(w http.ResponseWriter, r *http.Request) {
	tagid := mux.Vars(r)["tagid"]
	trabajos := []AtomPost{}
	err := db.Select(&trabajos, `SELECT publicaciones.id, publicaciones.fecha_publicacion, publicaciones.fecha_actualizacion, publicaciones.titulo, publicaciones.resumen, publicaciones.autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM publicaciones_tags LEFT JOIN publicaciones ON publicaciones_tags.publicacion_id = publicaciones.id LEFT JOIN autores ON publicaciones.autor_id = autores.id WHERE publicaciones_tags.tag_id = ? AND fecha_publicacion < NOW();`, tagid)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting trabajos: " + err.Error())
	}

	var tagnombre string
	err = db.QueryRowx(`SELECT nombre FROM tags WHERE id = ?;`, tagid).Scan(&tagnombre)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting trabajos: " + err.Error())
	}

	// TODO: Use tagid
	writeFeed(w, r, "tags-id-atom.xml", trabajos, tagnombre, tagid)
}

func AutorAtomFeed(w http.ResponseWriter, r *http.Request) {
	autorid := mux.Vars(r)["autorid"]

	var autor_nombre string
	err := db.QueryRowx("SELECT nombre FROM autores WHERE id=?", autorid).Scan(&autor_nombre)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Error("autoratom: feed not found")
		NotFoundHandler(w, r)
		return
	}

	tags := []Tag{}
	tagMap := map[string]string{}
	err = db.Select(&tags, `SELECT * FROM tags`)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting tags: " + err.Error())
	}

	for _, tag := range tags {
		tagMap[tag.Id] = tag.Nombre
	}

	posts := []AtomPost{}
	// TODO: Clean this query
	err = db.Select(&posts, `SELECT PublicacionesPublicas.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, tag_id, GROUP_CONCAT(tag_id) as raw_tags FROM PublicacionesPublicas LEFT JOIN publicaciones_tags ON PublicacionesPublicas.id = publicaciones_tags.publicacion_id LEFT JOIN autores ON PublicacionesPublicas.autor_id = autores.id WHERE autor_id=? GROUP BY id;`, autorid)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("unexpected error selecting posts by author: " + err.Error())
	}

	for i := 0; i < len(posts); i++ {
		p := posts[i]
		p.Tags = []string{}

		for _, tag := range strings.Split(p.Raw_tags.String, ",") {
			p.Tags = append(p.Tags, tagMap[tag])
		}

		posts[i] = p
	}

	writeFeed(w, r, "autores-atom.xml", posts, "Ariel Costas", autorid)
}

func writeFeed(w http.ResponseWriter, r *http.Request, feedName string, items []AtomPost, nombre string, id string) {
	// TODO: Refactor line above
	var lastUpdate time.Time

	for i := 0; i < len(items); i++ {
		p := &items[i]

		t, err := time.Parse("2006-01-02 15:04:05", p.Fecha_actualizacion)
		if err != nil {
			logger.Error("unexpected error parsing fecha_actualizacion: %s", err.Error())
			InternalServerErrorHandler(w, r)
		}

		if lastUpdate.Before(t) {
			lastUpdate = t
		}

		p.Id = url.PathEscape(p.Id)
	}

	w.Header().Add("Content-Type", "application/atom+xml; charset=utf-8")
	err := t.ExecuteTemplate(w, feedName, &FeedParams{
		BaseURL:    os.Getenv("DOMAIN"),
		LastUpdate: lastUpdate.Format(time.RFC3339),
		Posts:      items,
		Nombre:     nombre,
		Id:         id,
	})

	if err != nil {
		logger.Error("unexpected error rendering feed %s: %s", feedName, err.Error())
		InternalServerErrorHandler(w, r)
	}
}
