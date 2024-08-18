package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/httptesting"
)

func TestSignUp(t *testing.T) {
	signUpTestCases().Run(t)
}

func signUpTestCases() httptesting.TestingList {
	h := httptesting.Prepare(h)

	return httptesting.TestingList{
		h.Testing(httptesting.TestingSpecifications{
			Name: "Invalid email",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r, err := httptesting.NewRequestPostJson(pathSignUp, &app.Credentials{
					Email:    "invalid email",
					Password: "abcd1234",
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}
				return r
			},
			StatusCode: http.StatusBadRequest,
		}),
		h.Testing(httptesting.TestingSpecifications{
			Name: "Invalid password",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r, err := httptesting.NewRequestPostJson(pathSignUp, &app.Credentials{
					Email:    "w.matt.pavin@gmail.com",
					Password: "1",
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}
				return r
			},
			StatusCode: http.StatusBadRequest,
		}),
		h.Testing(httptesting.TestingSpecifications{
			Name: "Duplicated email",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				const dupEmail = "abcd1234@gmail.com"
				_ = signUp(t, dupEmail, "abcd1234")
				r, err := httptesting.NewRequestPostJson(pathSignUp, &app.Credentials{
					Email:    dupEmail,
					Password: "abcd1234",
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}
				return r
			},
			StatusCode: http.StatusBadRequest,
		}),
	}
}
