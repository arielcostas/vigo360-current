/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package templates

import (
	"bytes"
	"html/template"
	"strings"
	"time"
)

var functions = template.FuncMap{
	"safeHTML": func(text string) template.HTML {
		return template.HTML(text)
	},
	// Converts a standard date returned by MySQL to a RFC3339 datetime
	"date3339": func(date string) (string, error) {
		t, err := time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return "", err
		}
		return t.Format(time.RFC3339), nil
	},
	"dateDayMonth": func(date string) (string, error) {
		t, err := time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return "", err
		}
		return t.Format("02/01/2006"), nil
	},
	"markdown": func(text string) (template.HTML, error) {
		var buf bytes.Buffer
		err := parser.Convert([]byte(text), &buf)
		if err != nil {
			return template.HTML(""), err
		}
		return template.HTML(buf.Bytes()), nil
	},
	"split": func(text string, separator string) []string {
		return strings.Split(text, separator)
	},
	"iterateInt": func(num int) []int {
		var result []int
		for i := 1; i <= num; i++ {
			result = append(result, i)
		}
		return result
	},
	"wordCount": func(text string) int {
		return len(strings.Split(text, " "))
	},
}
