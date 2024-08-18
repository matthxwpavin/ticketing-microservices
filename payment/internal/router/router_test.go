package router

import (
	"context"
	"net/http"
	"os"
	"testing"

	mongodb "github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/payment/internal/database/mongo"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/matthxwpavin/ticketing/testsetup"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db *mongo.DB

var h http.Handler

var ctx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &mongo.DbConfig, func(loggerCtx context.Context, testdb *mongodb.DB) error {
		os.Setenv("STRIPE_SECRET", "sk_test_51PniCxFDYBNcFTuejDcU11Zn1qh1fAMpZqKW2lisFi4NoehleV1K0n104lTbYvR8jkQMJF5cR7OVTrbxEAb12tmc00dS8T8GXc")
		os.Setenv("JWT_KEY", "abcd")
		os.Setenv("NATS_URL", "nats://localhost:4222")
		os.Setenv("NATS_CONN_NAME", "some_name")
		if err := env.CheckRequiredEnvs([]env.EnvKey{
			env.StripeSecret,
			env.NatsURL,
			env.NatsConnName,
			env.JwtSecret,
		}); err != nil {
			return err
		}
		ctx = loggerCtx
		db = &mongo.DB{DB: testdb}
		h = New(ctx, db, &nats.MockClient{}, stripe.NewClient())
		return nil
	})
}

func addJwtCookieWithUserId(t *testing.T, r *http.Request, userId string) *http.Request {
	token, err := jwtclaims.IssueToken(jwtclaims.Metadata{
		UserID: userId,
		Email:  "abcd@xyz.com",
	})
	require.NoError(t, err, "could not generate a jwt token")

	r.AddCookie(jwtcookie.New(token))
	return r
}

func addJwtCookie(t *testing.T, r *http.Request) *http.Request {
	return addJwtCookieWithUserId(t, r, primitive.NewObjectID().Hex())
}
