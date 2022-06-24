/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

func NewMysqlTrabajoStore(db *sqlx.DB) *MysqlTrabajoStore {
	return &MysqlTrabajoStore{
		db: db,
	}
}

func (s *MysqlTrabajoStore) Listar() (models.Trabajos, error) {
	trabajos := make(models.Trabajos, 0)
	query := `SELECT t.id, COALESCE(fecha_publicacion, ""), fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM trabajos t LEFT JOIN autores ON t.autor_id = autores.id ORDER BY fecha_publicacion;`
	rows, err := s.db.Query(query)

	if err != nil {
		return trabajos, err
	}

	for rows.Next() {
		var nt models.Trabajo

		err = rows.Scan(&nt.Id, &nt.Fecha_publicacion, &nt.Fecha_actualizacion, &nt.Titulo, &nt.Resumen, &nt.Autor.Id, &nt.Autor.Nombre, &nt.Autor.Email)
		if err != nil {
			return models.Trabajos{}, err
		}

		trabajos = append(trabajos, nt)
	}
	return trabajos, nil
}

func (s *MysqlTrabajoStore) ListarPorAutor(autor_id string) (models.Trabajos, error) {
	var resultado = make(models.Trabajos, 0)
	trabajos, err := s.Listar()
	if err != nil {
		return models.Trabajos{}, err
	}

	for _, tr := range trabajos {
		if tr.Autor.Id == autor_id {
			resultado = append(resultado, tr)
		}
	}

	return resultado, nil
}

func (s *MysqlTrabajoStore) ObtenerPorId(id string, requirePublic bool) (models.Trabajo, error) {
	var post models.Trabajo
	var query = `SELECT trabajos.id, alt_portada, titulo, resumen, contenido, COALESCE(fecha_publicacion, ""), fecha_actualizacion, autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
	FROM trabajos
	LEFT JOIN autores on trabajos.autor_id = autores.id
	WHERE trabajos.id = ?
	GROUP BY trabajos.id 
	ORDER BY trabajos.fecha_publicacion DESC;`

	var err = s.db.QueryRow(query, id).Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion, &post.Autor.Id, &post.Autor.Nombre, &post.Autor.Biografia, &post.Autor.Rol)

	if err != nil {
		return models.Trabajo{}, err
	}

	if requirePublic && post.Fecha_publicacion == "" {
		return models.Trabajo{}, sql.ErrNoRows
	}

	return post, nil
}
