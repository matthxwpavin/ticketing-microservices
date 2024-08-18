package database

import "context"

type Database interface {
	OrderRepository() OrderRepository
	PaymentRepository() PaymentRepository
}

type OrderRepository interface {
	Insert(context.Context, *Order) (string, error)
	FindByID(context.Context, string) (*Order, error)
	UpdateByID(context.Context, string, *Order) error
}

type PaymentRepository interface {
	Insert(context.Context, *Payment) (string, error)
	FindByID(context.Context, string) (*Payment, error)
	FindByOrderId(context.Context, string) (*Payment, error)
	UpdateByID(context.Context, string, *Payment) error
	DeleteByID(context.Context, string) error
	FindByOrderIdAndStripePaymentIntentId(
		ctx context.Context,
		orderId string,
		stripePaymentIntentId string,
	) (*Payment, error)
	FindAll(context.Context) ([]*Payment, error)
}

type Order struct {
	Id      *string `bson:"_id,omitempty"`
	Version *int32  `bson:"version,omitempty"`
	Status  *string `bson:"status,omitempty"`
	UserId  *string `bson:"user_id,omitempty"`
	Price   *int32  `bson:"price,omitempty"`
}

type Payment struct {
	Id                    *string `bson:"_id,omitempty"`
	OrderId               *string `bson:"order_id,omitempty"`
	StripePaymentIntentId *string `bson:"stripe_payment_intent_id,omitempty"`
}
