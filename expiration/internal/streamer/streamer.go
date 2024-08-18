package streamer

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type Streamer interface {
	Disconenct(context.Context) error

	OrderCreatedConsumer(context.Context, streaming.ConsumeErrorHandler, string) (
		streaming.OrderCreatedConsumer,
		error,
	)
	ExpirationCompletedPublisher(context.Context) (streaming.ExpirationCompletedPublisher, error)
}
