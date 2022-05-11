package model

import "errors"

var ErrInvalidFormInput = errors.New("provided data is not valid")
var ErrExpiredSession = errors.New("session was older than 6 hours and was revoked automatically")
var ErrInvalidSession = errors.New("session token is not valid")
var ErrUnablePermissions = errors.New("unable to get permissions for session")
var ErrLoginRequired = errors.New("login is required to load this page")
var ErrBadParam = errors.New("param not provided or not valid")

type Session struct {
	Id           string
	Iniciada     string
	Autor_id     string
	Autor_nombre string
	Autor_rol    string
	Permisos     map[string]bool
}
