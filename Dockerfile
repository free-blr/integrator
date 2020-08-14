FROM golang:1.15-alpine3.12 AS builder

RUN apk update
RUN apk add --no-cache git make ca-certificates curl gcc g++

RUN go get -v github.com/rubenv/sql-migrate/sql-migrate

WORKDIR /src/free.blr/integrator
COPY . .

RUN PREFIX="/go/bin/" make build

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/sql-migrate /bin/sql-migrate
COPY --from=builder /src/free.blr/integrator/migrations/ /app/migrations/

COPY --from=builder /go/bin/bot /app/bot
