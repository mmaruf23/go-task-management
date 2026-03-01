package auth

import (
	"context"
	"errors"

	"github.com/mmaruf23/go-task-management/internal/db"
	"github.com/mmaruf23/go-task-management/pkg/util"
)

type AuthService struct {
	repo db.UserRepository
	jwt  JWTInterface
}

func NewAuthService(repo db.UserRepository, jwt JWTInterface) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (string, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return "", err
	}

	token, err := s.jwt.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.GenerateToken(user.ID.String())
}
