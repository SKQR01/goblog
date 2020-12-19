# goblog

## Database

### Database initial script

(In this case - Postgresql)

Postgres compare username with database name, so you need to login as "postgres" user with:

`$ sudo -i -u postgres` (after all operations you can return your user with "exite" command).

and create database named the same as your user(for example john, or foobar, or something else):

`$ psql`

`$ create_db skqr`

(in my case)

You can set up all DB manually or execute `setupDB.py` (on Linux).

### Migrate util setup

Download compiled program for your platfrom from here: https://github.com/golang-migrate/migrate/releases.

Then add (for Linux) this program to /usr/bin (you can rename it in way you want for your comfort).

For Windows you should to set PATH variable to your exe.

Creation of migrations:

`$ migrate create -ext sql -dir <path to migrations folder> <your migration name>`

Apply migrations:

`$ migrate -database postgres://<user>:<password>@<host:port>/<database> -path migrations up`(without quotes)

Drop migrations (in dirty migrate case):

`$ migrate -database postgres://<user>:<password>@<host:port>/<database> -path migrations force <migrate_version(you can know it from title of migration 20201216144859_init.up)> up(or down)`

You can find some development operations in deploy/setupDB.py file.



# Nginx setting up.

```bash
sudo apt-get install nginx nginx-common nginx-full
```


```bash
wget http://nginx.org/keys/nginx_signing.key
sudo apt-key add nginx_signing.key
```

```bash
sudo apt-get update
```

```bash
sudo apt-get install nginx
```

```bash
sudo systemctl start nginx 
```

```bash
#/etc/nginx/nginx.conf
events {}

http {
    include mime.types;

    server {
        listen 8080; #port of your future website
        location / { #proxy for... (to go to the site you need to go to the next url <your_ip_adress>:8080)
            proxy_pass http://127.0.0.1:8181; #url of your backend server (port must not matche with listen port)
        }
    }
}
```








