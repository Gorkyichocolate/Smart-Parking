FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/Smart-Parking/services/parking-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o parking-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/Smart-Parking/services/parking-service/parking-service .
EXPOSE 50053
CMD ["./parking-service"]
