package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/hash"
	tokenjwt "backend/pkg/token_jwt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: repo}
}

func (s *AuthService) SignUpClient(u *models.User, c echo.Context) (string, error) {
	count, err := s.authRepo.CheckUsersCount()
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("user already exists")
	}

	roleID, err := s.authRepo.GetAdminRoleId()
	if err != nil {
		return "", err
	}

	passwordHash, err := hash.GenerateHash(u.Password)
	if err != nil {
		return "", err
	}

	u.ID = uuid.New()
	u.RoleID = roleID
	u.PasswordHash = passwordHash

	token, err := tokenjwt.GenerateJWT(u.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, s.authRepo.CreateUser(u)
}

func (s *AuthService) SignInClient(u *models.User, c echo.Context) (models.AuthTokens, error) {
	client, err := s.authRepo.GetUserByEmailNumber(u.Email)
	if err != nil {
		return models.AuthTokens{}, err
	}
	if client == nil {
		return models.AuthTokens{}, errors.New("user not found")
	}
	if err := hash.ComparePassword(u.Password, client.PasswordHash); err != nil {
		return models.AuthTokens{}, errors.New("invalid password")
	}

	tokens, err := s.issueTokens(client.ID)
	if err != nil {
		return models.AuthTokens{}, err
	}

	*u = *client
	return tokens, nil
}

func (s *AuthService) RefreshToken(jwtStr string) (models.AuthTokens, error) {
	claims, err := tokenjwt.DecodeRefreshJWT(jwtStr)
	if err != nil {
		return models.AuthTokens{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return models.AuthTokens{}, fmt.Errorf("invalid user ID: %w", err)
	}
	refreshID, err := uuid.Parse(claims.RefreshId)
	if err != nil {
		return models.AuthTokens{}, fmt.Errorf("invalid refresh token ID: %w", err)
	}

	user, err := s.authRepo.GetUserById(userID)
	if err != nil {
		return models.AuthTokens{}, err
	}
	refresh, err := s.authRepo.GetRefreshTokenById(refreshID)
	if err != nil {
		return models.AuthTokens{}, err
	}

	if refresh.Secret != claims.Secret {
		return models.AuthTokens{}, errors.New("invalid refresh secret")
	}
	if refresh.ExpiresAt.Before(time.Now()) {
		return models.AuthTokens{}, errors.New("refresh token expired")
	}

	if err := s.authRepo.DeleteRefreshToken(refreshID); err != nil {
		return models.AuthTokens{}, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	return s.issueTokens(user.ID)
}

func (s *AuthService) SignOut(jwtStr string) error {
	claims, err := tokenjwt.DecodeRefreshJWT(jwtStr)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}
	refreshID, err := uuid.Parse(claims.RefreshId)
	if err != nil {
		return fmt.Errorf("invalid refresh token ID: %w", err)
	}
	refresh, err := s.authRepo.GetRefreshTokenById(refreshID)
	if err != nil {
		return fmt.Errorf("failed to get refresh token: %w", err)
	}
	if refresh == nil {
		return errors.New("refresh token not found")
	}
	if err := s.authRepo.DeleteRefreshToken(refreshID); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

// issueTokens — генерация и сохранение пары access/refresh токенов
func (s *AuthService) issueTokens(userID uuid.UUID) (models.AuthTokens, error) {
	accessToken, err := tokenjwt.GenerateJWT(userID)
	if err != nil {
		return models.AuthTokens{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := s.createRefreshToken(userID)
	if err != nil {
		return models.AuthTokens{}, err
	}

	return models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refresh,
	}, nil
}

// createRefreshToken — генерирует refresh token и сохраняет его в БД
func (s *AuthService) createRefreshToken(userID uuid.UUID) (string, error) {
	secret := GenerateSecret(12)
	refreshID := uuid.New()

	refreshToken, err := tokenjwt.GenerateRefreshJWT(userID, secret, refreshID)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refresh := &models.RefreshToken{
		ID:        refreshID,
		UserID:    userID,
		Secret:    secret,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.authRepo.SaveRefreshToken(refresh); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return refreshToken, nil
}

func GenerateSecret(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(uint64(time.Now().UnixNano()))

	secret := make([]byte, length)
	for i := range secret {
		secret[i] = charset[rand.Intn(len(charset))]
	}
	return string(secret)
}
