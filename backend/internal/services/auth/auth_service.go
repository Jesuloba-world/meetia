package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/meetia/backend/internal/models"
	"github.com/meetia/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenAuth *jwtauth.JWTAuth
	tokenExp  time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, secret string, tokenExp time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
		tokenExp:  tokenExp,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*models.User, string, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, "", ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// create new user
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExp).Unix(),
	}

	_, tokenString, err := s.tokenAuth.Encode(claims)
	return tokenString, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetTokenAuth() *jwtauth.JWTAuth {
	return s.tokenAuth
}
