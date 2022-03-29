## Instalación de Vigo360

1. Crear un nuevo usuario en el sistema y MySQL

```bash
sudo useradd -m -G www-data -s /bin/bash vigo360
sudo mysql
	> CREATE USER 'vigo360'@'localhost' IDENTIFIED BY '#P4s5w0rd';
	> CREATE SCHEMA vigo360 CHARACTER SET utf8mb4;
	> GRANT ALL PRIVILEGES ON vigo360.* TO USER 'vigo360';
	> exit;
sudo mkdir /var/www/vigo360 # O cualquier directorio, para los archivos subidos
sudo chown -R vigo360:www-data /var/www/vigo360
su - vigo360
```

2. Clonar y compilar

```bash
git clone https://gitlab.com/Vigo360/new.vigo360.es live
cd live
./launcher build
```

3. Copiar y activar el servicio systemd

```bash
sudo cp config/vigo360.service /etc/systemd/system/vigo360.service
sudo nano vigo360.service # O cualquier editor, modificar en [Service] los directorios correspondientes
```

4. Añadir configuración de NGINX

```bash
sudo cp config/nginx.conf /etc/nginx/sites-available/vigo360
sudo nano /etc/nginx/sites-available/vigo360 # Modificar dominio, directorios y demás
sudo ln -s /etc/nginx/sites-available/vigo360 /etc/nginx/sites-enabled/vigo360
sudo nginx -t
```

5. Ejecutar migraciones (pedirá contraseña de MySQL)

```bash
cat _schema/* | mysql -u vigo360 -p
```

6. Modificar variables de entorno

```bash
cp .env.example .env
nano .env # Modificar para que sea acorde a cada caso
```
