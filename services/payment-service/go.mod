module github.com/GorkyiChocolate/smart-parking/services/payment-service

go 1.21

require (
    github.com/GorkyiChocolate/smart-parking-proto/gen/go v0.0.0
    github.com/GorkyiChocolate/smart-parking/pkg v0.0.0
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    github.com/rabbitmq/amqp091-go v1.10.0
    google.golang.org/grpc v1.64.0
    google.golang.org/protobuf v1.34.1
)

replace (
    github.com/GorkyiChocolate/smart-parking-proto/gen/go => ../../../Smart-Parking-Proto/gen/go
    github.com/GorkyiChocolate/smart-parking/pkg => ../../pkg
)
