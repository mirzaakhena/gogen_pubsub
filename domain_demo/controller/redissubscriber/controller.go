package redissubscriber

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"log"
)

type controller struct {
	gogen.UsecaseRegisterer
	log logger.Logger
	cfg *config.Config

	funcHandlers map[string]func(msg *redis.Message)
}

func NewController(log logger.Logger, cfg *config.Config) gogen.ControllerRegisterer {

	return &controller{
		UsecaseRegisterer: gogen.NewBaseController(),
		log:               log,
		cfg:               cfg,
		funcHandlers:      map[string]func(msg *redis.Message){},
	}

}

func (r *controller) Start() {

	fmt.Println("REDIS PubSub Server is running on port 6379")

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	channels := make([]string, 0)
	for channelStr := range r.funcHandlers {
		channels = append(channels, channelStr)
	}

	pubsub := rdb.Subscribe(ctx, channels...)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Listen for messages.
	ch := pubsub.Channel()
	for msg := range ch {
		fh, ok := r.funcHandlers[msg.Channel]
		if !ok {
			continue
		}
		fh(msg)
	}

}
