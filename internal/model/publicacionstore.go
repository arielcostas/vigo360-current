/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import "github.com/jmoiron/sqlx"

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
	query := `SELECT p.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, GROUP_CONCAT(tags.nombre) as tags FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id GROUP BY id ORDER BY fecha_publicacion;`

	rows, err := s.db.Query(query)
	defer rows.Close()

	if err != nil {
		return publicaciones, err
	}

	for rows.Next() {
		var np Publicacion
		rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &np.Tags)
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

func (s *PublicacionStore) ObtenerPorId(id string) (Publicacion, error) {
	panic("not implemented")
}
