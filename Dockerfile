FROM golang:1.23-alpine

LABEL "org.opencontainers.image.source"="https://github.com/gianluparri03/cucinassistant"

WORKDIR /cucinassistant

COPY src/ ./

RUN go mod download
RUN go build main.go
RUN go build tools/broadcast.go

ENV CA_ENV=production
CMD ["./main"]
