package public

import (
	"database/sql"
	"errors"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

type AutoresParams struct {
	Autores []Autor
	Meta    common.PageMeta
}

type AutoresIdParams struct {
	Autor    Autor
	Posts    []ResumenPost
	Trabajos []ResumenPost
	Meta     common.PageMeta
}

func AutoresIdPage(w http.ResponseWriter, r *http.Request) {
	req_author := mux.Vars(r)["id"]
	autor := Autor{}

	err := db.QueryRowx("SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?", req_author).StructScan(&autor)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Error("[autores]: author not found with id %s", req_author)
		NotFoundHandler(w, r)
		return
	} else if err != nil {
		logger.Error("[autores]: unexpected error getting autor from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	publicaciones := []ResumenPost{}

	err = db.Select(&publicaciones, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM PublicacionesPublicas WHERE autor_id = ? ORDER BY PublicacionesPublicas.fecha_publicacion DESC;`, req_author)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[autores]: errors fetching posts from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	trabajos := []ResumenPost{}

	err = db.Select(&trabajos, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM TrabajosPublicos WHERE autor_id = ? ORDER BY TrabajosPublicos.fecha_publicacion DESC;`, req_author)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[autores]: errors fetching trabajos from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	t.ExecuteTemplate(w, "autores-id.html", AutoresIdParams{
		Autor:    autor,
		Posts:    publicaciones,
		Trabajos: trabajos,
		Meta: common.PageMeta{
			Titulo:      autor.Nombre,
			Descripcion: autor.Biografia,
			Canonica:    FullCanonica("/autores/" + autor.Id),
		},
	})
}

func AutoresPage(w http.ResponseWriter, r *http.Request) {
	autores := []Autor{}
	db.Select(&autores, `SELECT id, nombre, rol, biografia FROM autores;`)

	t.ExecuteTemplate(w, "autores.html", AutoresParams{
		Autores: autores,
		Meta: common.PageMeta{
			Titulo:      "Autores",
			Descripcion: "Conoce a los autores y colaboradores de Vigo360.",
			Canonica:    FullCanonica("/autores"),
		},
	})
}
