# Docker

This guide will explain to you how to setup an instance of CucinAssistant with
docker compose.

0. Make sure that both `docker` and `docker compose` are installed.

1. Create a config file (like `.env`) in your working directory, containing all
the required variables (see [configs/configs.go](../src/configs/configs.go) for
details).
```
CA_BASEURL="http://localhost:8080"
CA_PORT=8080
CA_SESSIONSECRET="random-string"
CA_DATABASE="user=ca password=ca dbname=ca host=database sslmode=disable"
CA_EMAIL_ENABLED=0
# ...
```

2. Then, create a `docker-compose.yml` file in the same directory, containing
```yaml
name: cucinassistant

services:
  app:
    build: .
    # or (if you want the latest stable version)
    # image: ghcr.io/gianluparri03/cucinassistant:latest

    depends_on:
      database:
        condition: service_healthy

    restart: always

    env_file: ".env"

    ports:
        - "127.0.0.1:8080:8080"

  database:
    image: postgres

    restart: always
    healthcheck:
      test: pg_isready

    volumes:
      - database:/var/lib/postgresql/data

    environment:
      POSTGRES_USER: ca
      POSTGRES_PASSWORD: ca
      POSTGRES_DB: ca

volumes:
  database:
```

3. Then,

**for the first-time install or after an update:**

a. Run `docker compose up -d database` to start the database

b. Run `docker compose run --rm app ca_migrate` to update/create the database schema

c. Run `docker compose up -d app` to start the app

**otherwise:**

Run `docker compose up -d`

4. It may happen that you need to tell something to your users. To do that, you
can simply execute `docker compose exec -it app ca_broadcast`. This 
will run a wizard that will ask you for the email subject and body, and then
(after a confirm) send it to everyone.
