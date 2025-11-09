package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"lalan-be/internal/config"
	"lalan-be/internal/response"
)

/*
Tipe untuk key context.
*/
type contextKey string

/*
Konstanta key untuk user ID di context.
*/
const UserIDKey contextKey = "user_id"

/*
Middleware untuk validasi JWT dan menambahkan user ID ke context.
*/
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Token required")
			return
		}

		// 2. Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Invalid token format")
			return
		}

		tokenString := parts[1]

		// 3. Parse & validasi token
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// 4. Simpan user ID ke context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.Subject)
		r = r.WithContext(ctx)

		// 5. Lanjut ke handler
		next.ServeHTTP(w, r)
	})
}

/*
Mengambil user ID dari context request.
*/
func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}
