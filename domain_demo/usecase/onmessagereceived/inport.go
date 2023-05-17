package onmessagereceived

import (
	"gogen_pubsub/shared/gogen"
)

type Inport = gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	Message string
}

type InportResponse struct {
}
