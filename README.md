# New Vigo360

New Vigo360 blog engine built in Go with MySQL.

## TODO

Things to implement before initial release.

### Admin

- [ ] Soporte para 2FA basada en TOTP
- [ ] Soporte para "series" de publicaciones

### Public

- [X] Soporte para OpenGraph
- [X] Sitemap.xml
- [ ] Feeds atom
	- [X] Publicaciones
	- [X] Por etiqueta
	- [ ] Por autor
	- [X] Trabajos
	- [ ] Links desde secciones por emoticono
- [ ] JSON-LD
	https://developers.google.com/search/docs/advanced/structured-data/search-gallery
	- [ ] Autor
	- [ ] Publicación
	- [ ] Inicio
- [X] Mostrar trabajos en perfil
- [ ] Errors: show requested path, function or something
- [ ] Errors: catch errors on page rendering and all SQL queries
- [ ] Soporte para "series" de publicaciones

- [ ] Post suggestions based on tags, author and date.
	- Point system? Score based on date, keywords, same author... Posts with most points win