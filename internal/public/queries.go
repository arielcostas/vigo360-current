/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

func GetFullTrabajo(id string) (Trabajo, error) {
	query := `SELECT trabajos.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
	FROM TrabajosPublicos trabajos 
	LEFT JOIN autores on trabajos.autor_id = autores.id 
	WHERE trabajos.id = ?;`

	trabajo := Trabajo{}
	if err := db.QueryRowx(query, id).StructScan(&trabajo); err != nil {
		return Trabajo{}, err
	}

	return trabajo, nil
}
