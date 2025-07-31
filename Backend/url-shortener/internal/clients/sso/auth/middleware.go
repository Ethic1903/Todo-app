package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type MiddlewareAuth struct {
	jwtSecret string
}

func New(jwtSecret string) *MiddlewareAuth {
	return &MiddlewareAuth{
		jwtSecret: jwtSecret,
	}
}

func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := a.extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		userID, err := a.validateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *MiddlewareAuth) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func (a *MiddlewareAuth) validateToken(token string) (int64, error) {
	checkToken, err := jwt.Parse(token, func(checkToken *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})
	if err != nil || !checkToken.Valid {
		return 0, err
	}

	claims := checkToken.Claims.(jwt.MapClaims)
	userID := int64(claims["uid"].(float64))
	return userID, nil
}
