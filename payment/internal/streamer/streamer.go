package streamer

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type Streamer interface {
	OrderCreatedConsumer(context.Context, streaming.ConsumeErrorHandler, string) (
		streaming.OrderCreatedConsumer,
		error,
	)
	OrderCancelledConsumer(context.Context, streaming.ConsumeErrorHandler, string) (
		streaming.OrderCancelledConsumer,
		error,
	)
	PaymentCreatedPublisher(context.Context) (streaming.PaymentCreatedPublisher, error)
}
