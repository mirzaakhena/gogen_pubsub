package rabbitmqsubscriber

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"log"
)

type controller struct {
	gogen.UsecaseRegisterer
	log logger.Logger
	cfg *config.Config

	funcHandlers map[string]func(msg *amqp091.Delivery)
}

func NewController(log logger.Logger, cfg *config.Config) gogen.ControllerRegisterer {

	return &controller{
		UsecaseRegisterer: gogen.NewBaseController(),
		log:               log,
		cfg:               cfg,
		funcHandlers:      map[string]func(msg *amqp091.Delivery){},
	}

}

func (r *controller) Start() {

	fmt.Println("RabbitMQ PubSub Server is running on port 9092")

	// Establish a connection to RabbitMQ server
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel on the connection
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare an exchange
	err = ch.ExchangeDeclare(
		"myExchange", // Exchange name
		"direct",     // Exchange type
		false,        // Durable
		false,        // Auto-delete
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	for routingKey := range r.funcHandlers {

		queue, err := ch.QueueDeclare(
			"",    // Generate a unique queue name
			false, // Durable
			false, // Auto-delete
			true,  // Exclusive
			false, // No-wait
			nil,   // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}

		err = ch.QueueBind(
			queue.Name,   // Queue name
			routingKey,   // Routing key
			"myExchange", // Exchange name
			false,        // No-wait
			nil,          // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to bind the queue: %v", err)
		}

		// Consume messages from the queue
		msgs, err := ch.Consume(
			queue.Name, // Queue name
			"",         // Consumer name
			true,       // Auto-acknowledge
			false,      // Exclusive
			false,      // No-local
			false,      // No-wait
			nil,        // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}

		// Start consuming messages for each routing key
		go func(routingKey string) {
			for msg := range msgs {

				fh, ok := r.funcHandlers[routingKey]
				if !ok {
					continue
				}

				fh(&msg)

			}
		}(routingKey)

	}

	// Wait indefinitely
	select {}

}
