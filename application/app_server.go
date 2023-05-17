package application

import (
	"gogen_pubsub/domain_demo/controller/kafkasubscriber"
	"gogen_pubsub/domain_demo/gateway/simpleprint"
	"gogen_pubsub/domain_demo/usecase/onmessagereceived"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
)

type appServer struct{}

func NewAppServer() gogen.Runner {
	return &appServer{}
}

func (appServer) Run() error {

	const appName = "appServer"

	cfg := config.ReadConfig()

	appData := gogen.NewApplicationData(appName)

	log := logger.NewSimpleJSONLogger(appData)

	datasource := simpleprint.NewGateway(log, appData, cfg)

	//primaryDriver := redissubscriber.NewController(log, cfg)
	primaryDriver := kafkasubscriber.NewController(log, cfg)

	primaryDriver.AddUsecase(
		//
		onmessagereceived.NewUsecase(datasource),
	)

	primaryDriver.RegisterRouter()

	primaryDriver.Start()

	return nil
}
