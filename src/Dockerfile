FROM golang:1.20-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o hroost-backend ./cmd

FROM ubuntu:22.04 AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y curl ca-certificates

COPY --from=builder /app/hroost-backend .

EXPOSE 3001

CMD ["./hroost-backend"]