/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
USE vigo360;

SET GLOBAL event_scheduler=ON;
DROP EVENT IF EXISTS expirar_sesiones;
CREATE EVENT expirar_sesiones
	ON SCHEDULE 
		EVERY 5 MINUTE 
	COMMENT 'Hace expirar las sesiones mayores de 6 horas automáticamente'
	DO
		UPDATE sesiones SET revocada = 0 WHERE iniciada < DATE_SUB(NOW(), INTERVAL 6 HOUR);