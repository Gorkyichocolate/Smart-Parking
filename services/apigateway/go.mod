module github.com/Gorkyichocolate/smart-parking/services/apigateway

go 1.21

require (
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth v0.0.0-00010101000000-000000000000
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/parking v0.0.0-00010101000000-000000000000
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/Gorkyichocolate/smart-parking-proto/gen/go/common v0.0.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go => ../../../Smart-Parking-Proto/gen/go

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth => ../../../Smart-Parking-Proto/gen/go/auth

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/parking => ../../../Smart-Parking-Proto/gen/go/parking

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment => ../../../Smart-Parking-Proto/gen/go/payment

replace github.com/Gorkyichocolate/smart-parking-proto/gen/go/common => ../../../Smart-Parking-Proto/gen/go/common
