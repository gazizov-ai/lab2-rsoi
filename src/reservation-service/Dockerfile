FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o reservation-service ./cmd/reservation-service

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/reservation-service /app/reservation-service
EXPOSE 8050

ENV PORT=8050

CMD ["/app/reservation-service"]
