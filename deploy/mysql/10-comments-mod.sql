USE vigo360;

CREATE VIEW vigo360.comment_moderation AS
SELECT c.id,
       c.publicacion_id,
       p.titulo                         as publicacion_titulo,
       COALESCE(padre_id, '')           as padre_id,
       c.nombre,
       c.es_autor,
       c.autor_original,
       c.contenido,
       c.fecha_creacion,
       COALESCE(c.fecha_moderacion, '') as fecha_moderacion,
       c.estado + 0                     as estado,
       COALESCE(c.moderador, '')        as moderador
FROM comentarios c
         LEFT JOIN publicaciones p ON c.publicacion_id = p.id;