FROM golang:1.23-alpine

LABEL "org.opencontainers.image.source"="https://github.com/gianluparri03/cucinassistant"

WORKDIR /cucinassistant

COPY go.mod go.sum ./
RUN go mod download

COPY configs/. configs/
COPY database/. database/
COPY email/. email/
COPY web/. web/
COPY main.go broadcast.go ./

RUN go build main.go
RUN go build broadcast.go

ENV CA_ENV=production
CMD ["./main"]
