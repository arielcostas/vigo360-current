/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"embed"
	"html/template"

	"vigo360.es/new/internal/logger"
)

//go:embed html/*
var rawtemplates embed.FS

var t = func() *template.Template {
	t := template.New("")

	functions := template.FuncMap{
		"safeHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
		"sum": func(a int, b int) int {
			return a + b
		},
	}

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
