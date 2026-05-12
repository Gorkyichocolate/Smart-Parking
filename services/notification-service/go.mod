module github.com/GorkyiChocolate/smart-parking/services/notification-service

go 1.21

require (
    github.com/GorkyiChocolate/smart-parking/pkg v0.0.0
    github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/GorkyiChocolate/smart-parking/pkg => ../../pkg
