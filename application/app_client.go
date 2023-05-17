package application

import (
	"gogen_pubsub/domain_demo/controller/restapi"
	"gogen_pubsub/domain_demo/gateway/kafkapublisher"
	"gogen_pubsub/domain_demo/usecase/runmessagesend"
	"gogen_pubsub/shared/config"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/infrastructure/token"
)

type appClient struct{}

func NewAppClient() gogen.Runner {
	return &appClient{}
}

func (appClient) Run() error {

	const appName = "appClient"

	cfg := config.ReadConfig()

	appData := gogen.NewApplicationData(appName)

	log := logger.NewSimpleJSONLogger(appData)

	jwtToken := token.NewJWTToken(cfg.JWTSecretKey)

	//datasource := redispublisher.NewGateway(log, appData, cfg)
	datasource := kafkapublisher.NewGateway(log, appData, cfg)

	primaryDriver := restapi.NewController(appData, log, cfg, jwtToken)

	primaryDriver.AddUsecase(
		//
		runmessagesend.NewUsecase(datasource),
	)

	primaryDriver.RegisterRouter()

	primaryDriver.Start()

	return nil
}
