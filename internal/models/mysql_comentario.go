/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

import "github.com/jmoiron/sqlx"

type MysqlComentarioStore struct {
	db *sqlx.DB
}

func NewMysqlComentarioStore(db *sqlx.DB) *MysqlComentarioStore {
	return &MysqlComentarioStore{
		db: db,
	}
}

// Lista los comentarios públicos para un artículo en forma de lista
func (s *MysqlComentarioStore) ListarPublicos(publicacion_id string) ([]Comentario, error) {
	var comentarios []Comentario
	var err = s.db.Select(&comentarios, `SELECT id, COALESCE(padre_id, "") as padre_id, nombre, es_autor, autor_original, contenido, fecha_moderacion FROM comentarios WHERE estado=? AND publicacion_id=?`, ESTADO_APROBADO, publicacion_id)
	if err != nil {
		return []Comentario{}, err
	}
	return comentarios, nil
}

// Lista los comentarios con un estado específico
func (s *MysqlComentarioStore) ListarPorEstado(estado EstadoComentario) ([]Comentario, error) {
	var comentarios []Comentario
	var err = s.db.Select(&comentarios, `SELECT * FROM comentarios WHERE estado=?`, estado)
	if err != nil {
		return []Comentario{}, err
	}
	return comentarios, nil
}

func (s *MysqlComentarioStore) GuardarComentario(_ Comentario) error {
	panic("not implemented") // TODO: Implement
}

// Cambia el estado de PENDIENTE a APROBADO
func (s *MysqlComentarioStore) Aprobar(comentario_id string) error {
	panic("not implemented") // TODO: Implement
}

// Cambia el estado de PENDIENTE a RECHAZADO
func (s *MysqlComentarioStore) Rechazar(comentario_id string) error {
	panic("not implemented") // TODO: Implement
}
