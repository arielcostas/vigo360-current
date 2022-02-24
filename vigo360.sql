-- 
-- WARNING: THIS SCRIPT WIPES THE DATABASE WITH ALL THE DATA
-- 
DROP SCHEMA IF EXISTS vigo360;
CREATE SCHEMA vigo360;
USE vigo360;

--
-- COMIENZA LA CREACION DE TABLAS
--
CREATE TABLE autores(
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

-- 
-- Rama publicaciones
-- 
CREATE TABLE publicaciones(
	id varchar(40) NOT NULL,
    fecha_publicacion datetime DEFAULT NULL,
    fecha_actualizacion datetime DEFAULT NOW() ON UPDATE NOW(),
    
    titulo varchar(80) NOT NULL,
    resumen varchar(200) NOT NULL,
    contenido text NOT NULL,

    autor_id varchar(40) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE palabras_clave(
	id int NOT NULL AUTO_INCREMENT,
    nombre varchar(40) NOT NULL,
    PRIMARY KEY (id)
);

-- 
-- Rama vehiculos
-- 
CREATE TABLE fotografias(
	id int NOT NULL AUTO_INCREMENT,
    titulo varchar(80) NOT NULL,
    descripcion varchar(500) NOT NULL,
    municipio varchar(40) NOT NULL,
    
    fecha_toma datetime NOT NULL,
    fecha_subida datetime NOT NULL,
    
    autor_id varchar(40) NOT NULL,
    vehiculo_id char(10) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE vehiculos(
	id char(10) NOT NULL,
    matricula char(10) NOT NULL UNIQUE,
    numeracion_empresa varchar(8) NOT NULL,
    comentario varchar(2000),
    
    tipo_vehiculo_id int NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE tipos_vehiculo(
	id int NOT NULL AUTO_INCREMENT,
	fabricante varchar(80) NOT NULL,
    modelo varchar(80) NOT NULL,
    titulo varchar(80) NOT NULL,
    comentario text,
    fecha_compra date NOT NULL,
    
    empresa_id int NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE empresas(
	id int NOT NULL AUTO_INCREMENT,
    nombre varchar(80) NOT NULL,
    descripcion varchar(2000) NOT NULL,
    municipio varchar(80) NOT NULL,
    
    PRIMARY KEY (id)
);

--
-- FOREIGN KEYS
--

-- usuario redacta publicacion
ALTER TABLE publicaciones ADD FOREIGN KEY publicaciones_autor(autor_id) REFERENCES autores(id);

-- publicacion pertenece a palabras clave
CREATE TABLE publicaciones_palabrasclave (
	publicacion_id varchar(40) NOT NULL,
    palabraclave_id int NOT NULL,
    PRIMARY KEY (publicacion_id, palabraclave_id),
    FOREIGN KEY ppc_publicacion (publicacion_id) REFERENCES publicaciones(id),
	FOREIGN KEY ppc_palabraclave(palabraclave_id) REFERENCES palabras_clave(id)
);

-- usuario comparte fotografia
ALTER TABLE fotografias ADD FOREIGN KEY fotografias_autor(autor_id) REFERENCES autores(id);

-- fotografia contiene vehiculo
ALTER TABLE fotografias ADD FOREIGN KEY fotografias_vehiculo(vehiculo_id) REFERENCES vehiculos(id);

-- vehiculo pertenece tipo_vehiculo
ALTER TABLE vehiculos ADD FOREIGN KEY vehiculos_tipovehiculo(tipo_vehiculo_id) REFERENCES tipos_vehiculo(id);

-- tipo_vehiculo propiedad de empresa
ALTER TABLE tipos_vehiculo ADD FOREIGN KEY tiposvehiculo_empresa(empresa_id) REFERENCES empresas(id);