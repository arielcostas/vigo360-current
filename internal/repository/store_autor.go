package repository

import "vigo360.es/new/internal/models"

type AutorStore interface {
	Listar() ([]models.Autor, error)
	Obtener(string) (models.Autor, error)
	Buscar(string) ([]models.Autor, error)
}
