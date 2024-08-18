package router

import (
	"context"
	"net/http"

	"github.com/matthxwpavin/ticketing/auth/internal/app"
	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/matthxwpavin/ticketing/serviceutil"
)

type handler struct {
	svc *app.Service
}

func (s *handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	creds := new(app.Credentials)
	ctx := r.Context()
	logger := sugar.FromContext(ctx)

	if err := rw.DecodeJSON(r.Body, creds); err != nil {
		logger.Errorw("Failed to read body", "error", err)
		rw.Error(ctx, w, serviceutil.NewServiceFailureError("Deocde JSON Error"))
		return
	}

	user, err := s.svc.SignUpUser(ctx, creds)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}

	if err := s.setJWTCookies(ctx, w, jwtclaims.Metadata{
		Email:  user.Email,
		UserID: user.ID,
	}); err != nil {
		return
	}

	rw.JSON201(ctx, w, user)

	logger.Infow("An user signin", "user ID", user.ID)
}

func (s *handler) handleSignin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := sugar.FromContext(ctx)

	creds := new(app.Credentials)
	if err := rw.DecodeJSON(r.Body, creds); err != nil {
		logger.Errorw("Failed to decode", "error", err)
		rw.Error(ctx, w, serviceutil.NewServiceFailureError("Deocde JSON Error"))
		return
	}

	user, err := s.svc.SignInUser(ctx, creds)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}

	if err := s.setJWTCookies(ctx, w, jwtclaims.Metadata{
		Email:  user.Email,
		UserID: user.ID,
	}); err != nil {
		return
	}

	rw.JSON(ctx, w, user)
}

func (s *handler) setJWTCookies(
	ctx context.Context,
	w http.ResponseWriter,
	metadata jwtclaims.Metadata,
) error {
	logger := sugar.FromContext(ctx)

	jwt, err := jwtclaims.IssueToken(metadata)
	if err != nil {
		logger.Errorw("Failed to sign token", "error", err)
		w.WriteHeader(500)
		return err
	}

	http.SetCookie(w, jwtcookie.New(jwt))
	return nil
}

func (s *handler) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := s.svc.CurrentUser(ctx)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON(ctx, w, user)
}

func (s *handler) handleSignOut(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, jwtcookie.New(""))
}
