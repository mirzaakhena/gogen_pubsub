package restapi

import (
	"context"
	"gogen_pubsub/domain_demo/usecase/runmessagesend"
	"gogen_pubsub/shared/gogen"
	"gogen_pubsub/shared/infrastructure/logger"
	"gogen_pubsub/shared/model/payload"
	"gogen_pubsub/shared/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *controller) runMessageSendHandler() gin.HandlerFunc {

	type InportRequest = runmessagesend.InportRequest
	type InportResponse = runmessagesend.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	type request struct {
		Message string `json:"message"`
	}

	type response struct {
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID(16)

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		err := c.BindJSON(&jsonReq)
		if err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var req InportRequest
		req.Message = jsonReq.Message

		r.log.Info(ctx, util.MustJSON(req))

		_, err = inport.Execute(ctx, req)
		if err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var jsonRes response

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes, traceID))

	}
}
