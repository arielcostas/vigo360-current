/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import (
	"time"

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
	var err = s.db.Select(&comentarios, `SELECT id, COALESCE(padre_id, "") as padre_id, nombre, es_autor, autor_original, contenido, COALESCE(fecha_moderacion, "") as fecha_moderacion FROM comentarios WHERE estado="aprobado" AND publicacion_id=? ORDER BY fecha_creacion ASC`, publicacion_id)
	if err != nil {
		return []models.Comentario{}, err
	}
	return comentarios, nil
}

// Lista los comentarios con un estado específico
func (s *MysqlComentarioStore) ListarPorEstado(estado models.EstadoComentario) ([]models.Comentario, error) {
	var comentarios []models.Comentario
	var err = s.db.Select(&comentarios, `SELECT id, publicacion_id, COALESCE(padre_id, "") as padre_id, nombre, es_autor, autor_original, contenido, fecha_creacion, COALESCE(fecha_moderacion, "") as fecha_moderacion, estado+0 as estado, COALESCE(moderador, "") as moderador FROM comentarios WHERE estado=?`, estado)
	if err != nil {
		return []models.Comentario{}, err
	}
	return comentarios, nil
}

func (s *MysqlComentarioStore) GuardarComentario(c models.Comentario) error {
	const query = `INSERT INTO comentarios(id, publicacion_id, padre_id, nombre, es_autor, autor_original, contenido, fecha_creacion, fecha_moderacion, estado, moderador) VALUES(?, ?, NULLIF(?, ""), ?, ?, ?,?,?,NULLIF(?, ''),?,NULLIF(?, ""))`
	if c.Fecha_creacion == "" {
		c.Fecha_creacion = time.Now().Format("2006-01-02 15:04:05")
	}
	_, err := s.db.Exec(query, c.Id, c.Publicacion_id, c.Padre_id, c.Nombre, c.Es_autor, c.Autor_original, c.Contenido, c.Fecha_creacion, c.Fecha_moderacion, c.Estado, c.Moderador)
	return err
}

// Cambia el estado de PENDIENTE a APROBADO
func (s *MysqlComentarioStore) Aprobar(comentario_id string, moderador string) error {
	_, err := s.db.Exec(`UPDATE comentarios SET estado=2, moderador=?, fecha_moderacion=NOW() WHERE id=? AND estado=1`, moderador, comentario_id)
	return err
}

// Cambia el estado de PENDIENTE a RECHAZADO
func (s *MysqlComentarioStore) Rechazar(comentario_id string, moderador string) error {
	_, err := s.db.Exec(`UPDATE comentarios SET estado=3, moderador=?,fecha_moderacion=NOW() WHERE id=? AND estado=1`, moderador, comentario_id)
	return err
}
