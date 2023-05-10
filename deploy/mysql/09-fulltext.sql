USE vigo360;

ALTER TABLE publicaciones
    DROP INDEX titulo;

ALTER TABLE publicaciones
    ADD FULLTEXT(id, titulo, resumen, contenido, alt_portada);
