package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/httptesting"
)

func TestSignOut(t *testing.T) {
	signOutTestCases().Run(t)
}

func signOutTestCases() httptesting.TestingList {
	ht := httptesting.Prepare(h)

	return httptesting.TestingList{
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Cookie is cleared",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				_ = signUp(t, "anique@one.com", "abcd1234")

				return httptesting.NewRequestPost(pathSignOut, nil)
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			if containsJWTCookie(r) {
				t.Errorf("JWT cookie is not cleared")
			}
		}),
	}

}
