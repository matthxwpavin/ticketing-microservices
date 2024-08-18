package stripe

import (
	"bytes"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/client"
	"github.com/stripe/stripe-go/v79/form"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

type Client struct {
	paymentintent *paymentintent.Client
}

func NewClient() *Client {
	return newClient(nil)
}

func NewClientTest() *Client {
	return newClient(&stripe.Backends{
		API:     &mockBackend{},
		Connect: &mockBackend{},
		Uploads: &mockBackend{},
	})
}

func newClient(backend *stripe.Backends) *Client {
	sc := &client.API{
		PaymentIntents: &paymentintent.Client{},
	}
	sc.Init(env.StripeSecret.Value(), backend)
	return &Client{
		paymentintent: sc.PaymentIntents,
	}
}

type mockBackend struct{}

func (m *mockBackend) Call(method, path, key string, params stripe.ParamsContainer, v stripe.LastResponseSetter) error {
	return nil
}

func (m *mockBackend) CallStreaming(method, path, key string, params stripe.ParamsContainer, v stripe.StreamingLastResponseSetter) error {
	return nil
}

func (m *mockBackend) CallRaw(method, path, key string, body *form.Values, params *stripe.Params, v stripe.LastResponseSetter) error {
	return nil
}

func (m *mockBackend) CallMultipart(method, path, key, boundary string, body *bytes.Buffer, params *stripe.Params, v stripe.LastResponseSetter) error {
	return nil
}

func (m *mockBackend) SetMaxNetworkRetries(maxNetworkRetries int64) {}
