package app

import (
	"context"
	"time"

	"github.com/matthxwpavin/ticketing/auth/internal/database"
	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/passwd"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is a type that represent domain/business logic
// of the application. It is high level language to communicate
// what exactly the application do.
type Service struct {
	ur database.UserRepository
}

func NewService(
	ur database.UserRepository,
) *Service {
	return &Service{
		ur: ur,
	}
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=14"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Service) SignUpUser(ctx context.Context, creds *Credentials) (*User, error) {
	logger := sugar.FromContext(ctx)
	if err := serviceutil.ValidateStruct(creds); err != nil {
		logger.Errorw("Credentials is invalid", "error", err)
		return nil, err
	}

	if user, err := s.ur.FindByEmail(ctx, creds.Email); err != nil {
		logger.Errorw("Failed to find an user", "error", err, "email", creds.Email)
		return nil, err
	} else if user != nil {
		logger.Errorw("Email in use", "email", creds.Email)
		return nil, serviceutil.NewServiceFailureError("Email in use")
	}

	hashedPasswd, err := passwd.Generate(creds.Password)
	if err != nil {
		logger.Errorw("Failed to generate a password", "error", err)
		return nil, err
	}

	user := &database.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     creds.Email,
		Password:  hashedPasswd,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if _, err := s.ur.Insert(ctx, user); err != nil {
		logger.Errorw("Failed to insert an user", "error", err)
		return nil, err
	}

	return &User{
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) SignInUser(ctx context.Context, creds *Credentials) (*User, error) {
	logger := sugar.FromContext(ctx)
	if err := serviceutil.ValidateStruct(creds); err != nil {
		logger.Errorw("Credentials is invalid", "error", err)
		return nil, err
	}

	user, err := s.ur.FindByEmail(ctx, creds.Email)
	if err != nil {
		logger.Errorw("Failed to find an user", "error", err)
		return nil, err
	}
	if user == nil {
		logger.Errorw("No user found", "email", creds.Email)
		return nil, serviceutil.NewServiceFailureError("No user found")
	}

	if err := passwd.Compare(user.Password, creds.Password); err != nil {
		logger.Errorw("Failed to compare password", "error", err)
		return nil, serviceutil.NewServiceFailureError("Failed to compare password")
	}

	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, nil
}

func (s *Service) CurrentUser(ctx context.Context) (*User, error) {
	logger := sugar.FromContext(ctx)
	claims := jwtclaims.FromContext(ctx)
	if claims == nil {
		logger.Errorw("No JWT claims found")
		return nil, serviceutil.ErrUnauthorized
	}

	return &User{
		ID:    claims.Metadata.UserID,
		Email: claims.Metadata.Email,
	}, nil
}
