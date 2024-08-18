package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/rw"
)

func TestCurrentUser(t *testing.T) {
	currentUserTestCases().Run(t)
}

func currentUserTestCases() httptesting.TestingList {
	ht := httptesting.Prepare(h)

	const email = "emailemail@emailemail.com"

	return httptesting.TestingList{
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Cookie current user",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				cookie := signUp(t, email, "abcd1234")

				r := httptest.NewRequest(http.MethodGet, pathCurrentUser, nil)
				r.AddCookie(cookie)
				return r
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			user := new(app.User)
			if err := rw.DecodeJSON(r.Body, user); err != nil {
				t.Errorf("failed to decode current user: %v", err)
			}
			if user.Email != email {
				t.Errorf("email is unexpected: expected: %v, received: %v", email, user.Email)
			}
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "No JWT cookie",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r := httptest.NewRequest(http.MethodGet, pathCurrentUser, nil)
				return r
			},
			StatusCode: http.StatusUnauthorized,
		}),
	}
}
