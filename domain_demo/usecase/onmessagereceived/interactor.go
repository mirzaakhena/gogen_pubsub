package onmessagereceived

import (
	"context"
)

type runMessageReverseInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runMessageReverseInteractor{
		outport: outputPort,
	}
}

func (r *runMessageReverseInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	r.outport.Print(ctx, req.Message)

	return res, nil
}
