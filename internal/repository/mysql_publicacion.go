package repository

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type MysqlPublicacionStore struct {
	db *sqlx.DB
}

func NewMysqlPublicacionStore(db *sqlx.DB) *MysqlPublicacionStore {
	return &MysqlPublicacionStore{
		db: db,
	}
}

func (s *MysqlPublicacionStore) Listar() (models.Publicaciones, error) {
	publicaciones := make(models.Publicaciones, 0)
	query := `SELECT p.id, COALESCE(fecha_publicacion, ""), fecha_actualizacion, titulo, resumen, alt_portada, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, COALESCE(GROUP_CONCAT(tags.id), "") as tags_ids, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags_nombres FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id GROUP BY id ORDER BY fecha_publicacion DESC;`

	rows, err := s.db.Query(query)

	if err != nil {
		return publicaciones, err
	}

	for rows.Next() {
		var (
			np            models.Publicacion
			rawTagIds     string
			rawTagNombres string
		)

		err = rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Alt_portada, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &rawTagIds, &rawTagNombres)
		if err != nil {
			return models.Publicaciones{}, err
		}

		var (
			tags            = make([]models.Tag, 0)
			splitTagIds     = strings.Split(rawTagIds, ",")
			splitTagNombres = strings.Split(rawTagNombres, ",")
		)
		for i := 0; i < len(splitTagIds); i++ {
			tags = append(tags, models.Tag{
				Id:     splitTagIds[i],
				Nombre: splitTagNombres[i],
			})
		}

		np.Tags = tags
		publicaciones = append(publicaciones, np)
	}
	return publicaciones, nil
}

func (s *MysqlPublicacionStore) ListarPorAutor(autor_id string) (models.Publicaciones, error) {
	var resultado = make(models.Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return models.Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		if pub.Autor.Id == autor_id {
			resultado = append(resultado, pub)
		}
	}

	return resultado, nil
}

func (s *MysqlPublicacionStore) ListarPorTag(tag_id string) (models.Publicaciones, error) {
	var resultado = make(models.Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return models.Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		for _, tag := range pub.Tags {
			if tag.Id == tag_id {
				resultado = append(resultado, pub)
				break
			}
		}
	}

	return resultado, nil
}

func (s *MysqlPublicacionStore) ListarPorSerie(serie_id string) (models.Publicaciones, error) {
	var resultado = make(models.Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return models.Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		if pub.Serie.Id == serie_id {
			resultado = append(resultado, pub)
		}
	}

	return resultado, nil
}

func (s *MysqlPublicacionStore) Existe(id string) (bool, error) {
	err := s.db.QueryRow("SELECT id FROM publicaciones WHERE id=? LIMIT 1", id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *MysqlPublicacionStore) ObtenerPorId(id string, requirePublic bool) (models.Publicacion, error) {
	var post models.Publicacion
	var query = `SELECT publicaciones.id, alt_portada, titulo, resumen, contenido, COALESCE(fecha_publicacion, ""), fecha_actualizacion, autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol, COALESCE(serie_id, ""), COALESCE(serie_posicion, 0) as serie_posicion, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags
	FROM publicaciones
	LEFT JOIN autores on publicaciones.autor_id = autores.id
	LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id
	LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id
	WHERE publicaciones.id = ?
	GROUP BY publicaciones.id 
	ORDER BY publicaciones.fecha_publicacion DESC;`

	var (
		rawTagNombres string
	)

	var err = s.db.QueryRow(query, id).Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion, &post.Autor.Id, &post.Autor.Nombre, &post.Autor.Biografia, &post.Autor.Rol, &post.Serie.Id, &post.Serie_posicion, &rawTagNombres)

	if err != nil {
		return models.Publicacion{}, err
	}

	if requirePublic && post.Fecha_publicacion == "" {
		return models.Publicacion{}, sql.ErrNoRows
	}

	for _, tag := range strings.Split(rawTagNombres, ",") {
		post.Tags = append(post.Tags, models.Tag{Nombre: tag})
	}

	return post, nil
}

func (s *MysqlPublicacionStore) Buscar(termino string) (models.Publicaciones, error) {
	var query = `SELECT p.id, COALESCE(fecha_publicacion, ""), fecha_actualizacion, titulo, resumen, alt_portada, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, COALESCE(GROUP_CONCAT(tags.id), "") as tags_ids, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags_nombres FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id WHERE MATCH(p.id, titulo, resumen, contenido, alt_portada) AGAINST (? IN NATURAL LANGUAGE MODE) GROUP BY id`

	rows, err := s.db.Query(query, "*"+termino+"*")
	if err != nil {
		return models.Publicaciones{}, err
	}

	var publicaciones = make(models.Publicaciones, 0)
	for rows.Next() {
		var (
			np            models.Publicacion
			rawTagIds     string
			rawTagNombres string
		)

		err = rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Alt_portada, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &rawTagIds, &rawTagNombres)
		if err != nil {
			return models.Publicaciones{}, err
		}

		var (
			tags            = make([]models.Tag, 0)
			splitTagIds     = strings.Split(rawTagIds, ",")
			splitTagNombres = strings.Split(rawTagNombres, ",")
		)
		for i := 0; i < len(splitTagIds); i++ {
			tags = append(tags, models.Tag{
				Id:     splitTagIds[i],
				Nombre: splitTagNombres[i],
			})
		}

		np.Tags = tags
		publicaciones = append(publicaciones, np)
	}

	return publicaciones, nil
}
