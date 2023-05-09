USE vigo360;

SET GLOBAL event_scheduler=ON;
DROP EVENT IF EXISTS expirar_sesiones;
CREATE EVENT expirar_sesiones
	ON SCHEDULE 
		EVERY 5 MINUTE 
	COMMENT 'Hace expirar las sesiones mayores de 6 horas automáticamente'
	DO
		UPDATE sesiones SET revocada = 0 WHERE iniciada < DATE_SUB(NOW(), INTERVAL 6 HOUR);