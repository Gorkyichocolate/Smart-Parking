module github.com/Gorkyichocolate/smart-parking/services/auth-service

go 1.21

require github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth v0.0.0

require (
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/common v0.0.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.33.0
)

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go => ../../../Smart-Parking-Proto/gen/go

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth => ../../../Smart-Parking-Proto/gen/go/auth

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/common => ../../../Smart-Parking-Proto/gen/go/common
