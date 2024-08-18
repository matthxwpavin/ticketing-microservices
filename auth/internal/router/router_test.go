package router

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/auth/internal/database/impl/mongo"
	mongodb "github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/testsetup"
)

var db *mongo.DB

var h http.Handler

var loggerCtx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, mongo.DbConfig, func(ctx context.Context, d *mongodb.DB) error {

		os.Setenv("JWT_KEY", "abcd")
		os.Setenv("NATS_URL", "nats://localhost:4222")
		os.Setenv("NATS_CONN_NAME", "some_name")
		os.Setenv("DEV", "dev")

		if err := env.CheckRequiredEnvs([]env.EnvKey{
			env.JwtSecret,
			env.NatsURL,
			env.NatsConnName,
			env.DEV,
		}); err != nil {
			return err
		}

		loggerCtx = ctx
		db = &mongo.DB{DB: d}
		h = New(loggerCtx, db)

		return nil
	})
}

func TestNotFound(t *testing.T) {
	httptesting.Run(t, httptesting.Testing{
		Handler: h,
		Specs: httptesting.TestingSpecifications{
			Name: "Not found",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return httptesting.NewRequestGet("/notfound-target")
			},
			StatusCode: http.StatusNotFound,
		},
	})
}

func signUp(t *testing.T, email, password string) *http.Cookie {
	var cookie *http.Cookie
	httptesting.Run(t, httptesting.Testing{
		Handler: h,
		Specs: httptesting.TestingSpecifications{
			Name: "SignUp helper",
			TestingRequest: func(t *testing.T) *http.Request {

				r, err := httptesting.NewRequestPostJson(pathSignUp, &app.Credentials{
					Email:    email,
					Password: password,
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}
				return r
			},
			StatusCode: http.StatusCreated,
		},
		AfterRun: func(t *testing.T, r *http.Response) {
			expectJWTCookie(t, r)
			cookie = findJWTCookie(r)
		},
	})
	return cookie
}

// func signUpDefault(t *testing.T) *http.Cookie {
// 	return signUp(t, "w.matt.pavin@gmail.com", "abcd1234")
// }

func expectJWTCookie(t *testing.T, r *http.Response) {
	if !containsJWTCookie(r) {
		t.Errorf("no JWT inside Set-Cookie header found")
	}
}

func containsJWTCookie(r *http.Response) bool {
	cookie := findJWTCookie(r)
	return cookie != nil && cookie.Value != ""
}

func findJWTCookie(r *http.Response) *http.Cookie {
	for _, cookie := range r.Cookies() {
		if cookie.Name == jwtcookie.Name {
			return cookie
		}
	}
	return nil
}
