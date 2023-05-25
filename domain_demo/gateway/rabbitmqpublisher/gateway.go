package kafkapublisher

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/model/payload"
)

type gateway struct {
	appData gogen.ApplicationData
	config  *config.Config
	log     logger.Logger
	channel *amqp091.Channel
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	// Establish a connection to RabbitMQ server
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a channel on the connection
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Declare an exchange
	err = ch.ExchangeDeclare(
		"myExchange", // Exchange name
		"fanout",     // Exchange type
		false,        // Durable
		false,        // Auto-delete
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		panic(err)
	}

	return &gateway{
		log:     log,
		appData: appData,
		config:  cfg,
		channel: ch,
	}
}

func (r *gateway) SendMessage(ctx context.Context, message string) error {

	r.log.Info(ctx, "called in Kafka Gateway")

	py := payload.Payload[payload.Message]{
		Data:      payload.Message{Content: message},
		Publisher: r.appData,
		TraceID:   logger.GetTraceID(ctx),
	}

	dataInBytes, err := json.Marshal(py)
	if err != nil {
		r.log.Error(ctx, err.Error())
		return err
	}

	topic := "sendMessage002"

	// Publish a message to the exchange
	err = r.channel.PublishWithContext(ctx,
		"myExchange", // Exchange name
		topic,        // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        dataInBytes,
		},
	)
	if err != nil {
		panic(err)
	}

	return nil

}
