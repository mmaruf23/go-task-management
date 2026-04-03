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

	userID, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", userID)
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

	userID, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", userID)
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

func TestGenerateToken_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	jti := uuid.New()
	mockRepo.On("CreateToken", mock.Anything, mock.Anything).Return(nil)

	jwtService := NewJWTService("iniceritanyasecret")
	service := NewAuthService(mockRepo, jwtService)

	userID := uuid.New()
	token, err := service.GenerateToken(ctx, jti, userID)

	assert.NoError(t, err)
	assert.NotEqual(t, "", token.Access)
	assert.NotEqual(t, "", token.Refresh)

	rtoken, err := jwtService.VerifyToken(token.Refresh)
	assert.NoError(t, err)
	assert.Equal(t, jti.String(), rtoken.ID)
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

func (m *MockUserRepo) CreateToken(ctx context.Context, arg repo.CreateTokenParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockUserRepo) ReplaceToken(ctx context.Context, arg repo.ReplaceTokenParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) RevokeToken(ctx context.Context, id uuid.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) RevokeAllToken(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
