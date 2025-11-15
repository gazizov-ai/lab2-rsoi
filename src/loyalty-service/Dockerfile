FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loyalty-service ./cmd/loyalty-service

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/loyalty-service /app/loyalty-service
EXPOSE 8070

ENV PORT=8070

CMD ["/app/loyalty-service"]
