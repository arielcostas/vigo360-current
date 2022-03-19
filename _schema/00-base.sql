CREATE SCHEMA IF NOT EXISTS vigo360_2;
USE vigo360_2;

CREATE TABLE IF NOT EXISTS autores(
	id varchar(40) NOT NULL,
    nombre varchar(40) NOT NULL,
    email varchar(150) NOT NULL UNIQUE,
    contraseña varchar(100) NOT NULL,

	rol varchar(40) NOT NULL,
    biografia varchar(2000) NOT NULL,
    web_url varchar(80) NOT NULL,
    web_titulo varchar(80) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS publicaciones(
	id varchar(40) NOT NULL,
    fecha_publicacion datetime DEFAULT NULL,
    fecha_actualizacion datetime DEFAULT NOW() ON UPDATE NOW(),
    alt_portada varchar(300) NOT NULL,
    
    titulo varchar(80) NOT NULL,
    resumen varchar(300) NOT NULL,
    contenido text NOT NULL,

    autor_id varchar(40) NOT NULL,
    serie_id varchar(40) DEFAULT NULL,
    serie_posicion int DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tags(
	id int NOT NULL AUTO_INCREMENT,
    nombre varchar(40) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS fotografias(
	id int NOT NULL AUTO_INCREMENT,
    titulo varchar(80) NOT NULL,
    descripcion varchar(500) NOT NULL,
    municipio varchar(40) NOT NULL,
    
    fecha_toma datetime NOT NULL,
    fecha_subida datetime NOT NULL,
    
    autor_id varchar(40) NOT NULL,
    vehiculo_id char(10) NOT NULL,
    licencia_id varchar(10) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS licencias(
	id varchar(10) NOT NULL,
    titulo varchar(50) NOT NULL,
    url varchar(100) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS vehiculos(
	id char(10) NOT NULL,
    matricula char(10) NOT NULL UNIQUE,
    numeracion_empresa varchar(8) NOT NULL,
    comentario varchar(2000),
    
    tipo_vehiculo_id int NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tipos_vehiculo(
	id int NOT NULL AUTO_INCREMENT,
	fabricante varchar(80) NOT NULL,
    modelo varchar(80) NOT NULL,
    titulo varchar(80) NOT NULL,
    comentario text,
    fecha_compra date NOT NULL,
    
    empresa_id int NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS empresas(
	id int NOT NULL AUTO_INCREMENT,
    nombre varchar(80) NOT NULL,
    descripcion varchar(2000) NOT NULL,
    municipio varchar(80) NOT NULL,
    
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS trabajos(
	id varchar(40) NOT NULL,
    titulo varchar(80) NOT NULL,
    resumen varchar(300) NOT NULL,
    contenido text NOT NULL,
    alt_portada varchar(300) NOT NULL,
    fecha_publicacion datetime NOT NULL DEFAULT NOW(),
    fecha_actualizacion datetime NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    
    autor_id varchar(40) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS adjuntos(
	id int NOT NULL AUTO_INCREMENT,
    nombre_archivo varchar(50) NOT NULL UNIQUE,
    titulo varchar(80) NOT NULL,
    fecha_subida datetime NOT NULL DEFAULT NOW(),
    trabajo_id varchar(40) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS sesiones(
    sessid char(20) NOT NULL,
    iniciada datetime NOT NULL DEFAULT NOW(),
    revocada boolean DEFAULT false,

    autor_id varchar(40) NOT NULL,
    PRIMARY KEY (sessid)
);

CREATE TABLE IF NOT EXISTS avisos(
    id int NOT NULL AUTO_INCREMENT,
    fecha_creacion datetime DEFAULT NOW(),
    titulo varchar(100) NOT NULL,
    contenido varchar(1000) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS series(
    id varchar(40) NOT NULL,
    titulo varchar(50) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS publicaciones_tags (
	publicacion_id varchar(40) NOT NULL,
    tag_id int NOT NULL,
    PRIMARY KEY (publicacion_id, tag_id),
    FOREIGN KEY ppc_publicacion (publicacion_id) REFERENCES publicaciones(id),
	FOREIGN KEY ppc_palabraclave(tag_id) REFERENCES tags(id)
);

-- autor redacta publicacion
ALTER TABLE publicaciones ADD FOREIGN KEY publicaciones_autor(autor_id) REFERENCES autores(id);

-- autor comparte fotografia
ALTER TABLE fotografias ADD FOREIGN KEY fotografias_autor(autor_id) REFERENCES autores(id);

-- fotografia contiene vehiculo
ALTER TABLE fotografias ADD FOREIGN KEY fotografias_vehiculo(vehiculo_id) REFERENCES vehiculos(id);

-- fotografia es cedida bajo licencia
ALTER TABLE fotografias ADD FOREIGN KEY fotografias_licencia(licencia_id) REFERENCES licencias(id);

-- vehiculo pertenece tipo_vehiculo
ALTER TABLE vehiculos ADD FOREIGN KEY vehiculos_tipovehiculo(tipo_vehiculo_id) REFERENCES tipos_vehiculo(id);

-- tipo_vehiculo propiedad de empresa
ALTER TABLE tipos_vehiculo ADD FOREIGN KEY tiposvehiculo_empresa(empresa_id) REFERENCES empresas(id);

-- trabajo contiene adjuntos
ALTER TABLE adjuntos ADD FOREIGN KEY adjuntos_trabajo(trabajo_id) REFERENCES trabajos(id);

-- autor inicia sesión
ALTER TABLE sesiones ADD FOREIGN KEY sesiones_autor(autor_id) REFERENCES autores(id);

-- publicación pertenece (opcionalmente) a serie
ALTER TABLE publicaciones ADD FOREIGN KEY publicaciones_series(serie_id) REFERENCES series(id);

CREATE VIEW sitemap AS SELECT uri, DATE(fecha_actualizacion) as fecha_actualizacion, priority, changefreq FROM (
	SELECT CONCAT("/tags/", publicaciones_tags.tag_id) as uri, fecha_actualizacion, "0.3" as priority, "weekly" as changefreq FROM publicaciones
		RIGHT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id
		WHERE fecha_actualizacion = (SELECT MAX(fecha_actualizacion) FROM publicaciones pub2
		RIGHT JOIN publicaciones_tags pt2 ON pub2.id = pt2.publicacion_id WHERE publicaciones_tags.tag_id = pt2.tag_id)
	UNION
	SELECT CONCAT("/autores/", publicaciones.autor_id) as uri, fecha_actualizacion, "0.7" as priority, "weekly" as changefreq FROM publicaciones
		WHERE fecha_actualizacion = (SELECT MAX(fecha_actualizacion) FROM publicaciones pub2
		WHERE publicaciones.autor_id = pub2.autor_id)
	UNION
	SELECT CONCAT("/trabajos/", trabajos.id), fecha_actualizacion, "0.5" as priority, "monthly" as changefreq FROM trabajos
	UNION
	SELECT CONCAT("/post/", publicaciones.id), fecha_actualizacion, "0.5" as priority, "monthly" as changefreq FROM publicaciones
	UNION
	(SELECT "/trabajos" as uri, fecha_actualizacion, "0.7" as priority, "weekly" as changefreq FROM trabajos ORDER BY fecha_actualizacion DESC LIMIT 1)
	UNION
	(SELECT "/" as uri, fecha_actualizacion, "0.9" as priority, "weekly" as changefreq FROM publicaciones ORDER BY fecha_actualizacion DESC LIMIT 1)
) as pq;