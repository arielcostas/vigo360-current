# New Vigo360

New Vigo360 blog engine built in Go with MySQL.

## TODO

Things to implement before initial release.

### Admin

-   [ ] Soporte para 2FA basada en TOTP
-   [x] Soporte para "series" de publicaciones
-   [-] Subida de imágenes extra a artículos

### Public

-   [ ] Feed atom por autor
-   [ ] Enlace a feeds atom desde HTML
-   [ ] Motor de búsqueda
    -   [ ] Soporte OpenSearch
-   [ ] Sistema de recomendación de publicaciones

### Otros

-   [ ] Errors: show requested path, function or something
-   [ ] Errors: catch errors on page rendering and all SQL queries
-   [ ] Cargar env desde systemd, no con godotenv

## Install guide

1. Clone via git and compile with `make build`
2. Add and enable systemd service (located in `config/`)
3. Add NGINX configuration
4. Run SQL migrations from `_schema/` "mysql -p < _schema/*"
5. Start systemd service and reload nginx config.
6. Check if it works.
