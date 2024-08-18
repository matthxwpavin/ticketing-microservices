package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/httptesting"
)

func TestSignIn(t *testing.T) {
	signInTestCases().Run(t)
}

func signInTestCases() httptesting.TestingList {
	h := httptesting.Prepare(h)

	return httptesting.TestingList{

		h.Testing(httptesting.TestingSpecifications{
			Name: "Email doesn't exists",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r, err := httptesting.NewRequestPostJson(pathSignIn, &app.Credentials{
					Email:    "somthing@notexists.com",
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

				const email = "valid.mail@yahoo.com"
				const fakePasswd = "a_fake_one"
				signUp(t, email, fakePasswd[1:])

				r, err := httptesting.NewRequestPostJson(pathSignIn, &app.Credentials{
					Email:    email,
					Password: fakePasswd,
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}
				return r

			},
			StatusCode: http.StatusBadRequest,
		}),

		h.Testing(httptesting.TestingSpecifications{
			Name: "Valid credentials given",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				const email = "valid2@abcd.com"
				const passwd = "abcd1234"
				_ = signUp(t, email, passwd)

				r, err := httptesting.NewRequestPostJson(pathSignIn, &app.Credentials{
					Email:    email,
					Password: passwd,
				})
				if err != nil {
					t.Errorf("failed to new request: %v", err)
				}

				return r
			},
			StatusCode: http.StatusOK,
		}).After(expectJWTCookie),
	}
}
