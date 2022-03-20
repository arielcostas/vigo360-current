# New Vigo360

New Vigo360 blog engine built in Go with MySQL.

## Install guide

1. Clone via git and compile with `make build`
2. Add and enable systemd service (located in `config/`)
3. Add NGINX configuration
4. Run SQL migrations from `_schema/` "mysql -p < _schema/*"
5. Start systemd service and reload nginx config.
6. Check if it works.
