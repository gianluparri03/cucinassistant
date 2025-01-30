# Docker

This guide will explain to you how to setup an instance of CucinAssistant with docker compose.

0. Make sure that both `docker` and `docker compose` are installed.

1. Create a config file (like `.env`) in your working directory, containing all the required
variables (see [configs/configs.go](../configs/configs.go) for details), like this:
```
CA_BASEURL="http://localhost:8080"
CA_PORT=80
CA_SESSIONSECRET="random-string"
CA_DATABASE="user=ca password=ca dbname=ca host=database sslmode=disable"
CA_EMAIL_ENABLED=0
# with CA_EMAIL_ENABLED=0, the server will not send any email. to enable it,
# you need to provide additional informations.
```

2. Then, create a `docker-compose.yml` file in the same directory, containing
```yaml
name: cucinassistant

services:
  app:
    build: .
    # or (if you want the latest stable version)
    # image: ghcr.io/gianluparri03/cucinassistant

    depends_on:
      database:
        condition: service_healthy

    restart: always
    healthcheck:
      test: "curl --fail localhost/info || exit 1"

    env_file: ".env"

    ports:
        - "127.0.0.1:8080:80"

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

3. Run `docker compose up` (with an optional `-d` to hide the output) and we're done!

4. It may happen that you need to tell something to your users. To do that, you can simply execute
`docker exec -it cucinassistant-app-1 ./broadcast config.yml`. This will run a wizard that will ask for the
email subject and body, then send a test email to the sender email. Then, after a confirm, it will broadcast it
to each user.
