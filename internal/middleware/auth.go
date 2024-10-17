package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/benpsk/go-survey-api/config"
	"github.com/benpsk/go-survey-api/pkg"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			pkg.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			pkg.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userId, ok := claims["user_id"].(float64)
		if !ok {
			pkg.Error(w, "Invalid user id", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
