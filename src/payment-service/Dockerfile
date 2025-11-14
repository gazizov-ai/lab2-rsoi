FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o payment-service ./cmd/payment-service

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/payment-service /app/payment-service
EXPOSE 8060

ENV PORT=8060

CMD ["/app/payment-service"]
