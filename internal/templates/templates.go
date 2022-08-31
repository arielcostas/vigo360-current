// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
)

//go:embed html/*
var rawtemplates embed.FS

var t = func() *template.Template {
	t := template.New("")

	entries, _ := rawtemplates.ReadDir("html")
	for _, de := range entries {
		filename := de.Name()
		contents, _ := rawtemplates.ReadFile("html/" + filename)

		_, err := t.New(filename).Funcs(Functions).Parse(string(contents))
		if err != nil {
			fmt.Printf("error parsing template: %s", err.Error())
			os.Exit(1)
		}
	}

	return t
}()

/*
	Render ejecuta una plantilla con los datos proveídos, llamando por debajo a ExecuteTemplate.
	Si hay un error al ejecutar la plantilla, no escribe nada al io.Writer y devuelve el error, con lo que es seguro no tener una página escrita a medias.
*/
func Render(w io.Writer, name string, data any) error {
	var output bytes.Buffer
	err := t.ExecuteTemplate(&output, name, data)
	if err != nil {
		return err
	}
	w.Write(output.Bytes())
	return nil
}
