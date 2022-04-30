/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"vigo360.es/new/internal/logger"
)

//go:embed html/*
var rawtemplates embed.FS

var t = func() *template.Template {
	t := template.New("")

	entries, _ := rawtemplates.ReadDir("html")
	for _, de := range entries {
		filename := de.Name()
		contents, _ := rawtemplates.ReadFile("html/" + filename)

		_, err := t.New(filename).Funcs(functions).Parse(string(contents))
		if err != nil {
			logger.Critical("[public-main] error parsing template: %s", err.Error())
		}
	}

	return t
}()

func Render(w io.Writer, name string, data any) error {
	var output bytes.Buffer
	err := t.ExecuteTemplate(&output, name, data)
	if err != nil {
		return err
	}
	w.Write(output.Bytes())
	return nil
}
