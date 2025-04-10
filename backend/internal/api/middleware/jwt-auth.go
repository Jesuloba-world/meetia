package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/jwtauth/v5"
)

func JWTMiddleware(tokenAuth *jwtauth.JWTAuth) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Unwrap only once at the beginning.
		r, w := humachi.Unwrap(ctx)

		// Get the standard Chi JWT middlewares
		jwtVerifier := jwtauth.Verifier(tokenAuth)
		jwtAuthenticator := jwtauth.Authenticator(tokenAuth)

		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx := huma.WithContext(ctx, r.Context())
			next(newCtx)
		})

		chainedHandler := jwtVerifier(jwtAuthenticator(finalHandler))
		chainedHandler.ServeHTTP(w, r)
	}
}
