# Instalación de Vigo360

1. Crear un nuevo usuario en el sistema y MySQL

```bash
sudo useradd -m -G www-data -s /bin/bash vigo360
sudo mysql
	> CREATE USER 'vigo360'@'localhost' IDENTIFIED BY '#P4s5w0rd';
	> CREATE SCHEMA vigo360 CHARACTER SET utf8mb4;
	> GRANT ALL PRIVILEGES ON vigo360.* TO USER 'vigo360';
	> exit;
cd /opt
sudo git clone https://github.com/arielcostas/vigo360.git vigo360
sudo chown -R vigo360:www-data vigo360/
su - vigo360
```

2. Clonar y compilar

```bash
cd /opt/vigo360
./launcher build
```

3. Copiar y activar el servicio systemd

```bash
sudo cp deploy/config/vigo360.service /etc/systemd/system/vigo360.service
sudo systemctl daemon-reload
sudo systemctl enable vigo360
```

4. Añadir configuración de NGINX

```bash
sudo cp config/nginx.conf /etc/nginx/sites-available/vigo360
sudo nano /etc/nginx/sites-available/vigo360 # Modificar dominio, ruta a certificados y puerto
sudo ln -s /etc/nginx/sites-available/vigo360 /etc/nginx/sites-enabled/vigo360
sudo nginx -t
```

5. Ejecutar migraciones (pedirá contraseña de MySQL)

```bash
cat deploy/mysql/* | mysql -D vigo360 -u vigo360 -p
```

6. Modificar variables de entorno

```bash
cp .env.example .env
openssl rand -hex 20 # Clave para indexnow, copiar y pegar en .env
nano .env # Modificar para que sea acorde a cada caso
```

7. Reiniciar nginx e iniciar servidor

```bash
sudo systemctl start vigo360.service
sudo nginx -s reload
```
