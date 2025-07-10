ARG GO_VERSION=1.24.1

FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["./main"]
