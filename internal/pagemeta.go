package internal

import "html/template"

type PageMeta struct {
	Titulo      string
	Descripcion string
	Keywords    string
	Canonica    string
	Miniatura   string
	BaseUrl     template.URL
}
