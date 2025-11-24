package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"lalan-be/internal/config"
	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
UserIDKey
konstanta kunci untuk user ID dalam konteks
*/
const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

/*
contextKey
type untuk kunci konteks dengan type safety
*/
type contextKey string

/*
Claims
struct untuk data JWT dengan role pengguna
*/
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/*
GetUserID
mengambil user ID dari konteks request
*/
func GetUserID(r *http.Request) string {
	if val := r.Context().Value(UserIDKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
GetUserRole
mengambil user role dari konteks request
*/
func GetUserRole(r *http.Request) string {
	if val := r.Context().Value(UserRoleKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
JWTMiddleware
memvalidasi token JWT dan memperbarui konteks dengan user data
*/
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("JWTMiddleware: Starting token validation for %s %s", r.Method, r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("JWTMiddleware: Authorization header missing")
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
		log.Printf("JWTMiddleware: Extracted token: %q", tokenString)
		if tokenString == "" || tokenString == authHeader {
			log.Printf("JWTMiddleware: Invalid token format")
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			log.Printf("JWTMiddleware: Token validation failed: %v", err)
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		log.Printf("JWTMiddleware: Token valid, userID: %s, role: %s", claims.Subject, claims.Role)
		log.Printf("JWTMiddleware: exp=%v (unix), now=%v", claims.ExpiresAt.Unix(), time.Now().Unix())
		log.Printf("JWTMiddleware: setting user_id to context: %s", claims.Subject)
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		r = r.WithContext(ctx)
		log.Printf("JWTMiddleware: context updated, calling next handler")

		next.ServeHTTP(w, r)
	})
}
