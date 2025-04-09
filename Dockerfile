FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -o auth_service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth_service .

COPY conf /app/conf
COPY migrations /app/migrations
COPY key /app/key

CMD ["./auth_service"]
