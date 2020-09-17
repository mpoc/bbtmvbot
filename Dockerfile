FROM golang:alpine AS builder
RUN apk add --no-cache build-base sqlite
WORKDIR /usr/src/bbtmvbot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cat createdb.sql | sqlite3 database.db
RUN go build -o bbtmvbot

FROM alpine
WORKDIR /app
COPY --from=builder /usr/src/bbtmvbot/bbtmvbot /usr/src/bbtmvbot/database.db /app/
ENTRYPOINT ["/app/bbtmvbot"]
