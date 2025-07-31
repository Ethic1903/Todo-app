package auth

import (
	"context"
	"net/http"
	"strings"
	ssogrpc "url-shortener/internal/clients/sso/grpc"
)

type AuthMiddleware struct {
	ssoClient *ssogrpc.Client
}

func New(ssoClient *ssogrpc.Client) *AuthMiddleware {
	return &AuthMiddleware{
		ssoClient: ssoClient,
	}
}

func (a *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := a.ssoClient.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
