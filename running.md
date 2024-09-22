# Running

This guide will explain to you how to setup an instance of CucinAssistant with docker compose.

0. Make sure that both `docker` and `docker compose` are installed.

1. Create a config file (like `prod.yml`) in your working directory, containing
```yaml
sessionSecret: "<a random string>"

baseURL: "<the host from which it will be accessible, like https://ca.gianlucaparri.me>"
port: "80"

database: "user=<db_user> password=<db_pass> dbname=<db_name> host=database sslmode=disable"

email:
  address: "<...>"
  server: "<...>"
  port: "<...>"
  password: "<...>"
```

2. Then, create a `docker-compose.yml` file in the same directory, like
```yaml
name: cucinassistant

services:
  app:
    image: github.com/gianluparri03/cucinassistant:latest

    depends_on:
      database:
        condition: service_healthy

    restart: always
    healthcheck:
      test: "curl --fail localhost/info || exit 1"

    volumes:
      - ./prod.yml:/cucinassistant/config.yml

    command: config.yml

    ports:
      - "<host_port>:80"

  database:
    image: postgres

    restart: always
    healthcheck:
      test: pg_isready

    environment:
      POSTGRES_USER: <db_user>
      POSTGRES_PASSWORD: <db_pass>
      POSTGRES_DB: <db_name>
```

3. Run `docker compose up` (with an optional `-d` to hide the output) and we're done!