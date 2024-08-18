package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

type ChargeInputI interface {
	Params() *stripe.PaymentIntentParams
}

type BaseChargeInput struct {
	Amount      *int64
	Description *string
}

func (i *BaseChargeInput) params(
	cur *string,
) *stripe.PaymentIntentParams {
	return &stripe.PaymentIntentParams{
		Amount:   i.Amount,
		Currency: cur,
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Description: i.Description,
	}
}

type ChargeInput struct {
	BaseChargeInput
	Currency *string
}

func (i *ChargeInput) Params() *stripe.PaymentIntentParams {
	return i.params(i.Currency)
}

type ChargeUsdInput struct {
	BaseChargeInput
}

func (i *ChargeUsdInput) Params() *stripe.PaymentIntentParams {
	return i.params(stripe.String(string(stripe.CurrencyUSD)))
}

func (c *Client) Charge(ctx context.Context, input ChargeInputI) (*stripe.PaymentIntent, error) {
	return c.paymentintent.New(input.Params())
}

func (c *Client) GetPaymentIntentById(ctx context.Context, id string) (*stripe.PaymentIntent, error) {
	return c.paymentintent.Get(id, &stripe.PaymentIntentParams{})
}
