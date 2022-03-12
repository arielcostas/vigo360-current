package admin

import (
	"regexp"
	"strings"
)

func validarId(id string) bool {
	return regexp.MustCompile(`^[a-z|ñ|\-|\_|0-9]{3,40}$`).MatchString(id)
}

func validarTitulo(titulo string) bool {
	return len(titulo) > 3 && len(titulo) < 80
}

func validarResumen(resumen string) bool {
	return len(resumen) > 3 && len(resumen) < 300
}

func validarContenido(contenido string) bool {
	return strings.TrimSpace(contenido) != ""
}
