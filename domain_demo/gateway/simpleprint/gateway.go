package simpleprint

import (
	"context"
	"fmt"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
)

type gateway struct {
	appData gogen.ApplicationData
	config  *config.Config
	log     logger.Logger
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	return &gateway{
		log:     log,
		appData: appData,
		config:  cfg,
	}
}

func (r gateway) Print(ctx context.Context, msg string) {
	fmt.Printf(">>> %s\n", msg)
}
