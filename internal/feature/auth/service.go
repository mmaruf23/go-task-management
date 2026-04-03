package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/mmaruf23/go-task-management/internal/repository"
	"github.com/mmaruf23/go-task-management/pkg/util"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error)
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error)
	UpdatePassword(ctx context.Context, arg repo.UpdatePasswordParams) (int64, error)
	CreateToken(ctx context.Context, arg repo.CreateTokenParams) error
	ReplaceToken(ctx context.Context, arg repo.ReplaceTokenParams) (int64, error)
	RevokeToken(ctx context.Context, id uuid.UUID) (int64, error)
	RevokeAllToken(ctx context.Context, userID uuid.UUID) (int64, error)
}

type AuthService struct {
	repo UserRepository
	jwt  *JWTService
}

type Token struct {
	Access              string
	Refresh             string
	ID                  uuid.UUID
	MaxAgeRefereshToken int
}

func NewAuthService(repo UserRepository, jwt *JWTService) *AuthService {
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

	return user.ID.String(), nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return user.ID.String(), nil
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

func (s *AuthService) GenerateToken(ctx context.Context, jti, userID uuid.UUID) (*Token, error) {
	maxAge := 14 * 24 * time.Hour
	exp := time.Now().Add(maxAge) // 2 minggu
	err := s.repo.CreateToken(ctx, repo.CreateTokenParams{
		ID:        jti,
		UserID:    userID,
		ExpiresAt: pgtype.Timestamp{Time: exp, Valid: true},
	})

	if err != nil {
		fmt.Printf("ERROR_CREATE_TOKEN : %s", err.Error())
		return nil, errors.New("auth process failed")
	}

	tokenID := jti.String()
	refreshToken, err := s.jwt.GenerateToken(userID.String(), exp, tokenID)
	if err != nil {
		fmt.Printf("ERROR_GENERATE_REFRESH_TOKEN : %s", err.Error())
		return nil, errors.New("auth process failed")
	}

	accessToken, err := s.jwt.GenerateToken(userID.String(), time.Now().Add(10*time.Minute), "")
	if err != nil {
		fmt.Printf("ERROR_GENERATE_ACCESS_TOKEN : %s", err.Error())
		return nil, errors.New("auth process failed")
	}

	token := Token{
		Access:              accessToken,
		Refresh:             refreshToken,
		ID:                  jti,
		MaxAgeRefereshToken: int(maxAge.Seconds()),
	}

	return &token, nil
}

func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	claims, err := s.jwt.VerifyToken(tokenString)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return errors.New("invalid token claims")
	}
	result, err := s.repo.RevokeToken(ctx, id)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("failed revoked token")
	}

	if result == 0 {
		return errors.New("no session was found.")
	}

	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {

	result, _ := s.repo.RevokeAllToken(ctx, userID)
	if result == 0 {
		return errors.New("no session was found.")
	}

	return nil
}

func (s *AuthService) ReplaceToken(ctx context.Context, tokenID, replacerID uuid.UUID) error {
	// todo next : tolong itu gimana caranya biar sqlc nggak generate pgtype.UUID buat replacerID nya. hadeh

	result, err := s.repo.ReplaceToken(ctx, repo.ReplaceTokenParams{
		ID:         tokenID,
		ReplacedBy: pgtype.UUID{Bytes: replacerID, Valid: true},
	})

	if err != nil {
		fmt.Printf("ERROR_REPLACE_TOKEN : %s", err.Error())
		return errors.New("auth failed")
	}

	if result == 0 {
		return errors.New("session missing / logged out")
	}
	return nil
}
