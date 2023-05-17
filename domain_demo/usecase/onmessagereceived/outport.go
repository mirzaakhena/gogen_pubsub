package onmessagereceived

import "context"

type Outport interface {
	Print(ctx context.Context, msg string)
}
