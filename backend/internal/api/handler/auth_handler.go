package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/meetia/backend/internal/models"
	"github.com/meetia/backend/internal/services/auth"
	humagroup "github.com/meetia/backend/lib/humaGroup"
)

type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(api huma.API) {
	authGroup := humagroup.NewHumaGroup(api, "/api/auth", []string{"Authentication"})

	humagroup.Post(authGroup, "/register", h.Register, "Register", &humagroup.HumaGroupOptions{
		Summary:     "Register a new user",
		Description: "Create a new user account with email, password, and display name",
	})
	humagroup.Post(authGroup, "/login", h.Login, "Login", &humagroup.HumaGroupOptions{
		Summary:     "User login",
		Description: "Authenticate a user with email and password",
	})
}

// RegisterBody defines the request body for user registration
type RegisterBody struct {
	Email       string `json:"email" doc:"User email address" example:"user@example.com"`
	Password    string `json:"password" doc:"User password" example:"securepassword123"`
	DisplayName string `json:"displayName" doc:"User display name" example:"John Doe"`
}

// LoginBody defines the request body for user login
type LoginBody struct {
	Email    string `json:"email" doc:"User email address" example:"user@example.com"`
	Password string `json:"password" doc:"User password" example:"securepassword123"`
}

// AuthResponse defines the response for authentication operations
type AuthResponse struct {
	User  *models.User `json:"user" doc:"User details"`
	Token string       `json:"token" doc:"JWT authentication token"`
}

type RegisterRequest struct {
	Body RegisterBody
}

type RegisterResponse struct {
	Body AuthResponse
}

func (h *AuthHandler) Register(ctx context.Context, input *RegisterRequest) (*RegisterResponse, error) {
	user, token, err := h.authService.Register(ctx, input.Body.Email, input.Body.Password, input.Body.DisplayName)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			return nil, huma.Error409Conflict("user already exists", err)
		default:
			return nil, huma.Error500InternalServerError("an error occured", err)
		}
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	return &RegisterResponse{
		Body: response,
	}, nil
}

type LoginRequest struct {
	Body LoginBody
}

type LoginResponse struct {
	Body AuthResponse
}

func (h *AuthHandler) Login(ctx context.Context, input *LoginRequest) (*LoginResponse, error) {
	user, token, err := h.authService.Login(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			return nil, huma.Error401Unauthorized("invalid credentials", err)
		default:
			return nil, huma.Error500InternalServerError("an error occured", err)
		}
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	return &LoginResponse{
		Body: response,
	}, nil
}
