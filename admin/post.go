package admin

import (
	"database/sql"
	"errors"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

type ResumenPost struct {
	Id           string
	Titulo       string
	Publicado    bool
	Autor_id     string
	Autor_nombre string
}

func PostListPage(w http.ResponseWriter, r *http.Request) {
	posts := []ResumenPost{}

	err := db.Select(&posts, `SELECT publicaciones.id, titulo, fecha_publicacion < NOW() as publicado, autor_id, autores.nombre as autor_nombre FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id ORDER BY fecha_publicacion DESC;`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("error inesperado leyendo publicaciones de la base de datos: %s", err.Error())
	}

	t.ExecuteTemplate(w, "post.html", struct {
		Posts []ResumenPost
	}{
		Posts: posts,
	})
}
