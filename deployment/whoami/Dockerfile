FROM golang:1.24.1-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build  -o whoami ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/whoami .

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chmod +x /app/whoami
USER appuser

CMD ["./whoami"]
