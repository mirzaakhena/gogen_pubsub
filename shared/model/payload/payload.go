package payload

import (
	"gogen_pubsub/shared/gogen"
)

type Payload[T any] struct {
	Data      T                     `json:"data"`
	Publisher gogen.ApplicationData `json:"publisher"`
	TraceID   string                `json:"traceId"`
}
