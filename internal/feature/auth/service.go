package auth

import (
	"context"
	"errors"

	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error)
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error)
	UpdatePassword(ctx context.Context, arg repo.UpdatePasswordParams) (int64, error)
}

type AuthService struct {
	repo UserRepository
	jwt  JWTInterface
}

func NewAuthService(repo UserRepository, jwt JWTInterface) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (string, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
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

func (s *AuthService) UpdatePassword(ctx context.Context, userID uuid.UUID, req *UpdatePasswordRequest) error {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	result, err := s.repo.UpdatePassword(ctx, repo.UpdatePasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})

	if err != nil {
		return err
	}

	if result == 0 {
		return errors.New("nothing was change") // just in case
	}

	return nil
}
