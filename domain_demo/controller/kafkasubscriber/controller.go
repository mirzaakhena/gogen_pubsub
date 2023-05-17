package kafkasubscriber

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
)

type controller struct {
	gogen.UsecaseRegisterer
	log logger.Logger
	cfg *config.Config

	funcHandlers map[string]func(msg *kafka.Message)
}

func NewController(log logger.Logger, cfg *config.Config) gogen.ControllerRegisterer {

	return &controller{
		UsecaseRegisterer: gogen.NewBaseController(),
		log:               log,
		cfg:               cfg,
		funcHandlers:      map[string]func(msg *kafka.Message){},
	}

}

func (r *controller) Start() {

	fmt.Println("Kafka PubSub Server is running on port 9092")

	// Kafka consumer configuration
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	})

	// close consumer later
	defer consumer.Close()

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}

	topics := make([]string, 0)
	for channelStr := range r.funcHandlers {
		topics = append(topics, channelStr)
	}

	// Subscribe to the topic
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return
	}

	// Start consuming messages
	for {

		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		fmt.Printf("Received message: %s\n", string(msg.Value))

		if msg.TopicPartition.Topic == nil {
			continue
		}

		fh, ok := r.funcHandlers[*msg.TopicPartition.Topic]
		if !ok {
			continue
		}

		fh(msg)

	}

}
