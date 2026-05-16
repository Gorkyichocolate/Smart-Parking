FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/Smart-Parking/services/apigateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o apigateway ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/Smart-Parking/services/apigateway/apigateway .
EXPOSE 8080
CMD ["./apigateway"]
