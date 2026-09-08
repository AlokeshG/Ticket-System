FROM golang:1.27-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o ticket-system .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

CMD ["./ticket-system"]