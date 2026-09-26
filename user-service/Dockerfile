FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/user-service ./cmd/app/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/bin/user-service .

USER appuser

EXPOSE 8080

CMD ["./user-service"]