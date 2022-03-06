package common

type NoPageData struct {
	Meta PageMeta
}

type PageMeta struct {
	Titulo      string
	Descripcion string
	Canonica    string
	Miniatura   string
}
