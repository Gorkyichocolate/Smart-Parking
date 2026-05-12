module github.com/GorkyiChocolate/smart-parking/services/payment-service

go 1.25.0

require (
	github.com/GorkyiChocolate/smart-parking-proto/gen/go v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.12.3
	github.com/rabbitmq/amqp091-go v1.11.0
	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
)

replace github.com/GorkyiChocolate/smart-parking-proto/gen/go => ../../../Smart-Parking-Proto/gen/go
