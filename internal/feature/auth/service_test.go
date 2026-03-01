package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"
	"github.com/stretchr/testify/assert"
)

// TEST
func TestRegister_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockUserRepository{
		createUserFunc: func(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
			return repo.User{
				ID:       uuid.New(),
				Name:     arg.Name,
				Email:    arg.Email,
				Password: arg.Password,
			}, nil
		},
	}

	jwtService := NewJWTService("secret")

	service := NewAuthService(mockRepo, jwtService)

	req := &RegisterRequest{
		Name:     "ini name test",
		Email:    "iniemail@test.com",
		Password: "inipasswordtest",
	}

	token, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", token)
}

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockUserRepository{
		getUserByEmailFunc: func(ctx context.Context, email string) (repo.User, error) {

			hashedPassword, err := util.HashPassword("password")
			if err != nil {
				return repo.User{}, err
			}

			return repo.User{
				ID:       uuid.New(),
				Name:     "user test",
				Email:    "user@test.com",
				Password: hashedPassword,
			}, nil
		},
	}

	jwtService := NewJWTService("iniceritanyasecret")

	service := NewAuthService(mockRepo, jwtService)

	req := &LoginRequest{
		Email:    "user@test.com",
		Password: "password",
	}

	token, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", token)
}

// MOCKING

type mockUserRepository struct {
	createUserFunc     func(ctx context.Context, arg repo.CreateUserParams) (repo.User, error)
	getUserByEmailFunc func(ctx context.Context, email string) (repo.User, error)
	getUserByIDFunc    func(ctx context.Context, id uuid.UUID) (repo.User, error)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	return m.createUserFunc(ctx, arg)
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return m.getUserByEmailFunc(ctx, email)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error) {
	return m.getUserByIDFunc(ctx, id)
}
