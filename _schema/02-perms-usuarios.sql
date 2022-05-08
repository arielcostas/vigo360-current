/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
USE vigo360;

CREATE TABLE permisos (
	id varchar(50) NOT NULL,
	comentario varchar(255) NOT NULL,
	PRIMARY KEY(id),
	CHECK(id != ""),
	CHECK(comentario != "")
);

CREATE TABLE permisos_usuarios (
	permiso_id varchar(50) NOT NULL,
	autor_id varchar(40) NOT NULL,
	PRIMARY KEY (permiso_id, autor_id),
	FOREIGN KEY (permiso_id) REFERENCES permisos(id),
	FOREIGN KEY (autor_id) REFERENCES autores(id)
);

/* Todos los permisos utilizables */
INSERT INTO permisos (id, comentario)
	VALUES ("publicaciones_delete", "Eliminar publicaciones");

/* La clave foránea no estaba puesta */
ALTER TABLE trabajos ADD FOREIGN KEY (autor_id) REFERENCES autores(id);
