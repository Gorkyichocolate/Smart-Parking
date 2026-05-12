// internal/queue/consumer.go
package queue

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumerConfig конфигурация потребителя
type ConsumerConfig struct {
	QueueName            string
	PrefetchCount        int
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
}

// Consumer управляет подключением к RabbitMQ
type Consumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	config         ConsumerConfig
	handler        MessageHandler
	stopChan       chan struct{}
	wg             sync.WaitGroup
	reconnectMutex sync.Mutex
	isConnected    bool
	rabbitMQURL    string // Добавляем поле для URL
}

// MessageHandler интерфейс обработчика сообщений
type MessageHandler interface {
	Handle(ctx context.Context, delivery amqp.Delivery) error
	GetQueueName() string
}

// NewConsumer создает нового потребителя
func NewConsumer(rabbitMQURL string, config ConsumerConfig, handler MessageHandler) (*Consumer, error) {
	consumer := &Consumer{
		config:      config,
		handler:     handler,
		stopChan:    make(chan struct{}),
		rabbitMQURL: rabbitMQURL,
	}

	if err := consumer.connect(rabbitMQURL); err != nil {
		return nil, err
	}

	return consumer, nil
}

// connect устанавливает соединение с RabbitMQ
func (c *Consumer) connect(rabbitMQURL string) error {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	// Объявляем очередь
	_, err = channel.QueueDeclare(
		c.config.QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	// Настройка QoS
	if c.config.PrefetchCount == 0 {
		c.config.PrefetchCount = 1
	}

	err = channel.Qos(
		c.config.PrefetchCount,
		0,
		false,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	c.conn = conn
	c.channel = channel
	c.isConnected = true

	log.Printf("✅ Connected to RabbitMQ, queue: %s", c.config.QueueName)
	return nil
}

// Start начинает обработку сообщений
func (c *Consumer) Start() error {
	c.wg.Add(1)
	go c.consume()
	return nil
}

// consume внутренний цикл потребления сообщений
func (c *Consumer) consume() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			log.Println("Stopping consumer...")
			return
		default:
			if !c.isConnected {
				if c.config.ReconnectDelay == 0 {
					c.config.ReconnectDelay = 5 * time.Second
				}
				time.Sleep(c.config.ReconnectDelay)

				// Пытаемся переподключиться
				if err := c.connect(c.rabbitMQURL); err != nil {
					log.Printf("Failed to reconnect: %v", err)
					continue
				}
			}

			// Начинаем потребление
			msgs, err := c.channel.Consume(
				c.config.QueueName,
				"",    // consumer tag
				false, // auto-ack
				false, // exclusive
				false, // no-local
				false, // no-wait
				nil,   // args
			)
			if err != nil {
				log.Printf("Failed to register consumer: %v", err)
				c.isConnected = false
				continue
			}

			// Обрабатываем сообщения
			for msg := range msgs {
				c.processMessage(msg)
			}

			// Если дошли сюда, канал закрыт
			c.isConnected = false
			log.Println("RabbitMQ channel closed, attempting to reconnect...")
		}
	}
}

// processMessage обрабатывает одно сообщение
func (c *Consumer) processMessage(msg amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Обрабатываем через handler
	err := c.handler.Handle(ctx, msg)

	if err != nil {
		log.Printf("❌ Failed to process message: %v", err)

		// Проверяем, можно ли повторить
		if msg.Redelivered {
			log.Printf("Message already redelivered, rejecting permanently")
			msg.Nack(false, false)
		} else {
			log.Printf("Requeue message for retry")
			msg.Nack(false, true)
		}
		return
	}

	// Успешная обработка
	msg.Ack(false)
}

// Stop останавливает потребителя
func (c *Consumer) Stop() {
	log.Println("Stopping consumer gracefully...")
	close(c.stopChan)
	c.wg.Wait()

	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}

	log.Println("Consumer stopped")
}
