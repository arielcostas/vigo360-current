/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type MysqlComentarioStore struct {
	db *sqlx.DB
}

func NewMysqlComentarioStore(db *sqlx.DB) *MysqlComentarioStore {
	return &MysqlComentarioStore{
		db: db,
	}
}

// Lista los comentarios públicos para un artículo en forma de lista
func (s *MysqlComentarioStore) ListarPublicos(publicacion_id string) ([]models.Comentario, error) {
	var comentarios []models.Comentario
	var err = s.db.Select(&comentarios, `SELECT id, COALESCE(padre_id, "") as padre_id, nombre, es_autor, autor_original, contenido, fecha_moderacion FROM comentarios WHERE estado="aprobado" AND publicacion_id=?`, publicacion_id)
	if err != nil {
		return []models.Comentario{}, err
	}
	return comentarios, nil
}

// Lista los comentarios con un estado específico
func (s *MysqlComentarioStore) ListarPorEstado(estado models.EstadoComentario) ([]models.Comentario, error) {
	var comentarios []models.Comentario
	var err = s.db.Select(&comentarios, `SELECT * FROM comentarios WHERE estado=?`, estado)
	if err != nil {
		return []models.Comentario{}, err
	}
	return comentarios, nil
}

func (s *MysqlComentarioStore) GuardarComentario(c models.Comentario) error {
	const query = `INSERT INTO comentarios VALUES(?, ?, ?, ?, ?, ?,?,?,?,?,?)`
	_, err := s.db.Exec(query, c.Id, c.Publicacion_id, c.Padre_id, c.Nombre, c.Es_autor, c.Autor_original, c.Contenido, c.Fecha_creacion, c.Fecha_moderacion, c.Estado, c.Moderador)
	return err
}

// Cambia el estado de PENDIENTE a APROBADO
func (s *MysqlComentarioStore) Aprobar(comentario_id string) error {
	_, err := s.db.Exec(`UPDATE comentarios SET estado=1 WHERE comentario_id=?`, comentario_id)
	return err
}

// Cambia el estado de PENDIENTE a RECHAZADO
func (s *MysqlComentarioStore) Rechazar(comentario_id string) error {
	_, err := s.db.Exec(`UPDATE comentarios SET estado=2 WHERE comentario_id=?`, comentario_id)
	return err
}
