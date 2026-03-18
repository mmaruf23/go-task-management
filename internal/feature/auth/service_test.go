package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TEST
func TestRegister_Success(t *testing.T) {
	ctx := context.Background()

	params := repo.CreateUserParams{
		Name:     "ini name test",
		Email:    "iniemail@test.com",
		Password: "inipasswordtest",
	}

	mockRepo := new(MockUserRepo)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(
		repo.User{
			ID:       uuid.New(),
			Name:     params.Name,
			Email:    params.Email,
			Password: params.Password,
		}, nil)

	jwtService := NewJWTService("secret")

	service := NewAuthService(mockRepo, jwtService)

	req := &RegisterRequest{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}

	token, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", token)
}

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockUserRepo)
	hashedPassword, _ := util.HashPassword("password")
	mockRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repo.User{
		ID:       uuid.New(),
		Name:     "user test",
		Email:    "user@test.com",
		Password: hashedPassword,
	}, nil)

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

func TestUpdatePassword_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockRepo.On("UpdatePassword", mock.Anything, mock.Anything).Return(int64(1), nil)
	jwtService := NewJWTService("iniceritanyasecret")

	service := NewAuthService(mockRepo, jwtService)
	req := &UpdatePasswordRequest{
		Password: "iniupdatedpassword",
	}

	userID := uuid.New()
	err := service.UpdatePassword(ctx, userID, req)

	assert.NoError(t, err)
}

// MOCKING

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.User), args.Error(1)
}
func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(repo.User), args.Error(1)
}
func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repo.User), args.Error(1)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, arg repo.UpdatePasswordParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}
