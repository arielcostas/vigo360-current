package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type TagStore interface {
	Listar() ([]models.Tag, error)
	Obtener(string) (models.Tag, error)
}

type MysqlTagStore struct {
	db *sqlx.DB
}

func NewMysqlTagStore(db *sqlx.DB) *MysqlTagStore {
	return &MysqlTagStore{
		db: db,
	}
}

func (s *MysqlTagStore) Listar() ([]models.Tag, error) {
	var tags = make(map[string]models.Tag, 0)
	var rows, err = s.db.Query(`SELECT id, nombre FROM tags`)
	if err != nil {
		return []models.Tag{}, err
	}

	for rows.Next() {
		var nt models.Tag
		err = rows.Scan(&nt.Id, &nt.Nombre)
		if err != nil {
			return []models.Tag{}, err
		}
		tags[nt.Id] = nt
	}

	rows, err = s.db.Query(`SELECT tag_id, COUNT(publicacion_id) FROM publicaciones_tags GROUP BY tag_id`)
	if err != nil {
		return []models.Tag{}, err
	}
	for rows.Next() {
		var t string
		var c int

		rows.Scan(&t, &c)
		var nt = tags[t]
		nt.Publicaciones = c
		tags[t] = nt
	}

	var tagSlice []models.Tag
	for _, t := range tags {
		tagSlice = append(tagSlice, t)
	}
	return tagSlice, nil
}

func (s *MysqlTagStore) Obtener(tag_id string) (models.Tag, error) {
	var tag models.Tag
	var row = s.db.QueryRow(`SELECT id, nombre FROM tags WHERE id=?`, tag_id)
	var err = row.Scan(&tag.Id, &tag.Nombre)
	if err != nil {
		return models.Tag{}, err
	}

	row = s.db.QueryRow(`SELECT COUNT(*) FROM publicaciones_tags WHERE tag_id=?`, tag_id)
	err = row.Scan(&tag.Publicaciones)
	if err != nil {
		return models.Tag{}, err
	}

	return tag, nil
}
