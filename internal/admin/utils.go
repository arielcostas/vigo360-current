/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"vigo360.es/new/internal/logger"
)

func ComprobarContraseña(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}

	if errors.Is(err, bcrypt.ErrHashTooShort) {
		logger.Notice("[validatepassword]: unable to verify password: hash is too short")
	}

	return false
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
