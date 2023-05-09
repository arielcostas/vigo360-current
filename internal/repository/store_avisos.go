package repository

import "vigo360.es/new/internal/models"

type AvisoStore interface {
	// Obtiene todos los avisos
	Listar() ([]models.Aviso, error)
	// Obtiene los 5 avisos más recientes
	ListarRecientes() ([]models.Aviso, error)
}
