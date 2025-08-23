FROM golang:1.24.1-alpine AS builder

RUN go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
FROM alpine:latest

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

RUN migrate -version

ENTRYPOINT [ "/usr/local/bin/migrate" ]
