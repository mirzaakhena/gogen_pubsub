package redispublisher

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/model/payload"
)

type gateway struct {
	appData gogen.ApplicationData
	config  *config.Config
	log     logger.Logger
	client  *redis.Client
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &gateway{
		log:     log,
		appData: appData,
		config:  cfg,
		client:  client,
	}
}

func (r *gateway) SendMessage(ctx context.Context, message string) error {

	r.log.Info(ctx, "called in Redis Gateway")

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

	err = r.client.Publish(ctx, "sendMessage001", string(dataInBytes)).Err()
	if err != nil {
		r.log.Error(ctx, err.Error())
		return err
	}

	return nil

}
