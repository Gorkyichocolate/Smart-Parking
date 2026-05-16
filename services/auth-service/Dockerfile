FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/Smart-Parking/services/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/Smart-Parking/services/auth-service/auth-service .
EXPOSE 50051
CMD ["./auth-service"]
