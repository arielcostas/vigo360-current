package repository

import "vigo360.es/new/internal/models"

type SerieStore interface {
	Listar() ([]models.Serie, error)
	Obtener(string) (models.Serie, error)
}
