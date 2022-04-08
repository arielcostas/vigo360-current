/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
	"math/rand"
	"strings"
)

/*
 Generates three suggested posts based on certain criteria, like matching author, same tags...
 For more details, check https://gitlab.com/vigo360/new.vigo360.es/-/issues/2
*/
type PostRecommendation struct {
	ResumenPost
	Tags   sql.NullString
	Points int
}

// Compares both slices for matching tags, and returns how many match
func FindMatchingTags(tags1, tags2 []string) int {
	// Make it a map to avoid passing it many times
	var tag1map = make(map[string]bool)
	for _, v := range tags1 {
		tag1map[v] = true
	}

	// Iterate tag2 to find which match
	var matches = 0

	for _, v := range tags2 {
		if _, ok := tag1map[v]; ok {
			matches++
		}
	}

	return matches
}

func generateSuggestions(original_id string) ([]PostRecommendation, error) {
	var resultado = make([]PostRecommendation, 3)

	// Two by points
	var original PostRecommendation
	err := db.QueryRowx(`SELECT pp.id, fecha_publicacion, titulo, autores.id as autor_id, GROUP_CONCAT(tag_id) as tags FROM PublicacionesPublicas pp LEFT JOIN autores ON pp.autor_id = autores.id LEFT JOIN publicaciones_tags ON pp.id = publicaciones_tags.publicacion_id WHERE pp.id = ? GROUP BY publicacion_id;`, original_id).StructScan(&original)
	var original_tags = strings.Split(original.Tags.String, ",")
	if err != nil {
		return []PostRecommendation{}, err
	}

	// Rest
	var options []PostRecommendation
	err = db.Select(&options, `SELECT pp.id, DATE_FORMAT(fecha_publicacion, '%d %b.') as fecha_publicacion, titulo, autores.id as autor_id, resumen, alt_portada, autores.nombre, GROUP_CONCAT(tag_id) as tags FROM PublicacionesPublicas pp LEFT JOIN autores ON pp.autor_id = autores.id LEFT JOIN publicaciones_tags ON pp.id = publicaciones_tags.publicacion_id WHERE pp.id != ? GROUP BY pp.id;`, original_id)
	if err != nil {
		return []PostRecommendation{}, err
	}

	var tags int
	err = db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&tags)

	for i, rp := range options {
		var points = 10

		/*	Same author => +12 points
			Different => 0 */
		if original.Autor_id == rp.Autor_id {
			points += 12
		}

		/*	If all tags match, give 3 times as many points tags are
		 *	Also, for each that matches give 2 points
		 *	3 tags all match => +9 points
		 *	2 match => +4 points
		 */
		matches := FindMatchingTags(original_tags, strings.Split(rp.Tags.String, ","))
		if len(original_tags) == matches {
			points += len(original_tags) * 4
		}

		points += matches * 3

		/*	Some random points to not make them the same all the time
			Adds or deduces up to 8 points, randomly */
		points += rand.Intn(12) - 6

		// Persist it
		rp.Points = points
		options[i] = rp
	}

	resultado[0] = options[0]
	for _, rp := range options[1:] {
		if resultado[0].Points < rp.Points {
			resultado[1] = resultado[0]
			resultado[0] = rp
			continue
		}

		if resultado[1].Points < rp.Points {
			resultado[1] = rp
		}
	}

	// Random suggestion
	err = db.QueryRowx(`SELECT pp.id, DATE_FORMAT(fecha_publicacion, '%d %b.') as fecha_publicacion, alt_portada, titulo, resumen, nombre FROM PublicacionesPublicas pp  LEFT JOIN autores ON pp.autor_id = autores.id WHERE pp.id != ? ORDER BY RAND() LIMIT 1;`, original_id).StructScan(&resultado[2])
	if err != nil {
		return []PostRecommendation{}, err
	}
	return resultado, nil
}
