package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/jwtauth/v5"

)

type AuthParam struct {
	Authorization string `required:"true" header:"Authorization" example:"Bearer <token>" doc:"Bearer <token>"`
}

func getUserIdFromContext(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", huma.Error401Unauthorized("invalid token", err)
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", huma.Error401Unauthorized("invalid token claims")
	}

	return userID, nil
}
