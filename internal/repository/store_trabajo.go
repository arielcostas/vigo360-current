package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type TrabajoStore interface {
	Listar() (models.Trabajos, error)
	ListarPorAutor(string) (models.Trabajos, error)
	ObtenerPorId(string, bool) (models.Trabajo, error)
}

type MysqlTrabajoStore struct {
	db *sqlx.DB
}
