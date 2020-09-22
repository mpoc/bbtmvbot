FROM golang:alpine AS builder
RUN apk add --no-cache build-base sqlite
WORKDIR /usr/src/bbtmvbot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bbtmvbot
RUN cat createdb.sql | sqlite3 database.db

FROM alpine
WORKDIR /app
COPY --from=builder /usr/src/bbtmvbot/bbtmvbot /usr/src/bbtmvbot/database.db /app/
ENTRYPOINT ["/app/bbtmvbot"]
