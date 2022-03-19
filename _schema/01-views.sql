CREATE VIEW PublicacionesPublicas AS
SELECT * FROM publicaciones
WHERE publicaciones.fecha_publicacion IS NOT NULL AND publicaciones.fecha_publicacion < NOW();

CREATE VIEW TrabajosPublicos AS
SELECT * FROM trabajos
WHERE trabajos.fecha_publicacion IS NOT NULL AND trabajos.fecha_publicacion < NOW();
