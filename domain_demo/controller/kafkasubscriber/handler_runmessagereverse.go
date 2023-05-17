package kafkasubscriber

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gogen_pubsub/domain_demo/usecase/onmessagereceived"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/model/payload"
	"gogen_pubsub/shared/util"
)

func (r *controller) sendMessageHandler(msg *kafka.Message) {

	type InportRequest = onmessagereceived.InportRequest
	type InportResponse = onmessagereceived.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	traceID := util.GenerateID(16)

	ctx := logger.SetTraceID(context.Background(), traceID)

	obj := payload.Payload[payload.Message]{}
	err := json.Unmarshal(msg.Value, &obj)
	if err != nil {
		r.log.Error(ctx, err.Error())
		return
	}

	var req InportRequest
	req.Message = obj.Data.Content

	_, err = inport.Execute(ctx, req)
	if err != nil {
		r.log.Error(ctx, err.Error())
		return
	}

}
