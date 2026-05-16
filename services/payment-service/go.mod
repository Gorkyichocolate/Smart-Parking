module github.com/Gorkyichocolate/smart-parking/services/payment-service

go 1.21

replace github.com/Gorkyichocolate/smart-parking/pkg => ../../pkg

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go => ../../../Smart-Parking-Proto/gen/go

replace github.com/Gorkyichocolate/smart-parking/pkg/config => ../../pkg/config

replace github.com/Gorkyichocolate/smart-parking/pkg/metrics => ../../pkg/metrics

require github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment v0.0.0

require (
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/common v0.0.0 // indirect
	github.com/Gorkyichocolate/smart-parking/pkg/config v0.0.0-20260512100515-3b77ecdae3f9
	github.com/Gorkyichocolate/smart-parking/pkg/metrics v0.0.0-20260512100515-3b77ecdae3f9
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.33.0
)

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment => ../../../Smart-Parking-Proto/gen/go/payment

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/common => ../../../Smart-Parking-Proto/gen/go/common
