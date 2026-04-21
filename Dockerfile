FROM golang:1.26-alpine AS build

WORKDIR /cucinassistant
COPY src/ ./

RUN apk add make
RUN make install
RUN make gen

RUN go build main.go
RUN go build tools/broadcast.go
RUN go build tools/migrate.go


FROM alpine:latest AS run

RUN apk add curl

COPY --from=build /cucinassistant/main /bin/cucinassistant
COPY --from=build /cucinassistant/broadcast /bin/ca_broadcast
COPY --from=build /cucinassistant/migrate /bin/ca_migrate

HEALTHCHECK CMD curl --fail localhost/info || exit 1

ENV CA_ENV="production"

CMD ["cucinassistant"]
