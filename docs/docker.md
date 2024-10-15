# Docker

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
    image: ghcr.io/gianluparri03/cucinassistant:latest

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

    volumes:
      - database:/var/lib/postgresql/data

    environment:
      POSTGRES_USER: <db_user>
      POSTGRES_PASSWORD: <db_pass>
      POSTGRES_DB: <db_name>

volumes:
  database:
```

3. Run `docker compose up` (with an optional `-d` to hide the output) and we're done!

4. It may happen that you need to tell something to your users. To do that, you can simply execute
`docker exec -it cucinassistant-app-1 ./broadcast config.yml`. This will run a wizard that will ask for the
email subject and body, then send a test email to the sender email. Then, after a confirm, it will broadcast it
to each user.
