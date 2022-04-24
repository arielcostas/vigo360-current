USE vigo360;

ALTER TABLE publicaciones ADD FULLTEXT(titulo, resumen, contenido);
