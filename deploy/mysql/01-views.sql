/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
USE vigo360;

CREATE VIEW PublicacionesPublicas AS
SELECT * FROM publicaciones
WHERE publicaciones.fecha_publicacion IS NOT NULL AND publicaciones.fecha_publicacion < NOW();

CREATE VIEW TrabajosPublicos AS
SELECT * FROM trabajos
WHERE trabajos.fecha_publicacion IS NOT NULL AND trabajos.fecha_publicacion < NOW();
