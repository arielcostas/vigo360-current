/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PublicacionStore struct {
	db *sqlx.DB
}

func NewPublicacionStore(db *sqlx.DB) PublicacionStore {
	return PublicacionStore{
		db: db,
	}
}

func (s *PublicacionStore) Listar() (Publicaciones, error) {
	publicaciones := make(Publicaciones, 0)
	query := `SELECT p.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, GROUP_CONCAT(tags.id) as tags_ids, GROUP_CONCAT(tags.nombre) as tags_nombres FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id GROUP BY id ORDER BY fecha_publicacion;`

	rows, err := s.db.Query(query)

	if err != nil {
		return publicaciones, err
	}

	for rows.Next() {
		var (
			np            Publicacion
			rawTagIds     string
			rawTagNombres string
		)

		err = rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &rawTagIds, &rawTagNombres)
		if err != nil {
			return Publicaciones{}, err
		}

		var (
			tags            = make([]Tag, 0)
			splitTagIds     = strings.Split(rawTagIds, ",")
			splitTagNombres = strings.Split(rawTagNombres, ",")
		)
		for i := 0; i < len(splitTagIds); i++ {
			tags = append(tags, Tag{
				Id:     splitTagIds[i],
				Nombre: splitTagNombres[i],
			})
		}

		np.Tags = tags
		publicaciones = append(publicaciones, np)
	}
	return publicaciones, nil
}

func (s *PublicacionStore) ListarPorAutor(autor_id string) (Publicaciones, error) {
	var resultado = make(Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		if pub.Autor.Id == autor_id {
			resultado = append(resultado, pub)
		}
	}

	return resultado, nil
}

func (s *PublicacionStore) ListarPorTag(tag_id string) (Publicaciones, error) {
	var resultado = make(Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return Publicaciones{}, err
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

func (s *PublicacionStore) ListarPorSerie(serie_id string) (Publicaciones, error) {
	var resultado = make(Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		if pub.Serie.Id == serie_id {
			resultado = append(resultado, pub)
		}
	}

	return resultado, nil
}

func (s *PublicacionStore) ObtenerPorId(id string, requirePublic bool) (Publicacion, error) {
	var post Publicacion
	var query = `SELECT publicaciones.id, alt_portada, titulo, resumen, contenido, COALESCE(fecha_publicacion, ""), fecha_actualizacion, autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol, COALESCE(serie_id, ""), GROUP_CONCAT(tags.nombre) as tags
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

	var err = s.db.QueryRow(query, id).Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion, &post.Autor.Id, &post.Autor.Nombre, &post.Autor.Biografia, &post.Autor.Rol, &post.Serie.Id, &rawTagNombres)

	if err != nil {
		return Publicacion{}, err
	}

	if requirePublic && post.Fecha_publicacion == "" {
		return Publicacion{}, sql.ErrNoRows
	}

	for _, tag := range strings.Split(rawTagNombres, ",") {
		post.Tags = append(post.Tags, Tag{Nombre: tag})
	}

	return post, nil
}
