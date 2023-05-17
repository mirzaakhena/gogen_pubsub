package kafkapublisher

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/model/payload"
)

type gateway struct {
	appData  gogen.ApplicationData
	config   *config.Config
	log      logger.Logger
	producer *kafka.Producer
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	return &gateway{
		log:      log,
		appData:  appData,
		config:   cfg,
		producer: producer,
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
	err = r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          dataInBytes,
	}, nil)

	if err != nil {
		r.log.Error(ctx, err.Error())
		return err
	}

	return nil

}
