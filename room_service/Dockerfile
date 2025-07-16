FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -o room_service ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/room_service .

COPY conf /app/conf
COPY migrations /app/migrations
COPY key /app/key

CMD ["./room_service"]