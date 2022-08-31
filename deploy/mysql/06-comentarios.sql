USE vigo360;

ALTER TABLE comentarios DROP COLUMN email_hash, MODIFY COLUMN moderador varchar(40) DEFAULT NULL;