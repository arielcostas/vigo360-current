package public

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
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
	// TODO error handling
	err := db.QueryRowx("SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?", req_author).StructScan(&autor)

	if errors.Is(err, sql.ErrNoRows) {
		NotFoundHandler(w, r)
		return
	} else if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	publicaciones := []ResumenPost{}
	// TODO error handling
	err = db.Select(&publicaciones, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM publicaciones WHERE autor_id = ? ORDER BY publicaciones.fecha_publicacion DESC;`, req_author)

	if err != nil {
		log.Fatalf(err.Error())
	}

	trabajos := []ResumenPost{}
	// TODO error handling
	err = db.Select(&trabajos, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM trabajos WHERE autor_id = ? ORDER BY trabajos.fecha_publicacion DESC;`, req_author)

	if err != nil {
		log.Fatalf(err.Error())
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
