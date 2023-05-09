package models

type Web struct {
	Url    string
	Titulo string
}

type Autor struct {
	Id        string
	Nombre    string
	Email     string
	Rol       string
	Biografia string
	Web       Web

	Publicaciones Publicaciones
}
