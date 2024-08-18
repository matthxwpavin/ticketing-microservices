package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/auth/internal/database"
	"github.com/matthxwpavin/ticketing/middleware"
)

const (
	group = "/api/users"

	pathCurrentUser = group + "/currentuser"
	pathSignUp      = group + "/signup"
	pathSignIn      = group + "/signin"
	pathSignOut     = group + "/signout"
)

func New(ctx context.Context, db database.Database) *mux.Router {
	r := mux.NewRouter()

	handler := &handler{
		svc: app.NewService(db.UserRepository()),
	}

	r.Use(middleware.PopulateLogger)
	r.Use(middleware.PopulateJWTClaims)

	r.HandleFunc(pathCurrentUser, handler.handleCurrentUser).Methods(http.MethodGet)
	r.HandleFunc(pathSignUp, handler.handleSignup).Methods(http.MethodPost)
	r.HandleFunc(pathSignIn, handler.handleSignin).Methods(http.MethodPost)
	r.HandleFunc(pathSignOut, handler.handleSignOut).Methods(http.MethodPost)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	return r
}
