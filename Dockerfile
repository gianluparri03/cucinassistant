FROM golang:1.23-alpine AS build

WORKDIR /cucinassistant
COPY src/ ./

RUN go mod download
RUN go build main.go
RUN go build tools/broadcast.go
RUN go build tools/migrate.go


FROM alpine:latest

RUN apk --no-cache add curl

COPY --from=build /cucinassistant/main /bin/cucinassistant
COPY --from=build /cucinassistant/broadcast /bin/ca_broadcast
COPY --from=build /cucinassistant/migrate /bin/ca_migrate

HEALTHCHECK CMD curl --fail localhost/info || exit 1

ENV CA_ENV="production"
CMD ["cucinassistant"]
