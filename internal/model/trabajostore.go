/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import (
	"github.com/jmoiron/sqlx"
)

type TrabajoStore struct {
	db *sqlx.DB
}

func NewTrabajoStore(db *sqlx.DB) TrabajoStore {
	return TrabajoStore{
		db: db,
	}
}

func (s *TrabajoStore) Listar() (Trabajos, error) {
	trabajos := make(Trabajos, 0)
	query := `SELECT t.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM trabajos t LEFT JOIN autores ON t.autor_id = autores.id ORDER BY fecha_publicacion;`
	rows, err := s.db.Query(query)

	if err != nil {
		return trabajos, err
	}

	for rows.Next() {
		var nt Trabajo

		rows.Scan(&nt.Id, &nt.Fecha_publicacion, &nt.Fecha_actualizacion, &nt.Titulo, &nt.Resumen, &nt.Autor.Id, &nt.Autor.Nombre, &nt.Autor.Email)

		trabajos = append(trabajos, nt)
	}
	return trabajos, nil
}

func (s *TrabajoStore) ListarPorAutor(autor_id string) (Trabajos, error) {
	panic("por implementar")
}

func (s *TrabajoStore) ObtenerPorId() (Trabajos, error) {
	panic("por implementar")
}
