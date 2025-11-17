package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"lalan-be/internal/config"
	"lalan-be/internal/response"
)

/*
Konstanta untuk kunci konteks.
Mendefinisikan kunci user ID dan role dalam konteks.
*/
const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

/*
contextKey digunakan sebagai kunci untuk nilai konteks.
Memastikan type safety dalam context.
*/
type contextKey string

/*
Claims berisi data JWT standar dan role pengguna.
Digunakan untuk parsing token JWT.
*/
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/*
GetUserID mengambil user ID dari konteks.
Mengembalikan string kosong jika tidak ada.
*/
func GetUserID(r *http.Request) string {
	if val := r.Context().Value(UserIDKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
GetUserRole mengambil user role dari konteks.
Mengembalikan string kosong jika tidak ada.
*/
func GetUserRole(r *http.Request) string {
	if val := r.Context().Value(UserRoleKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
JWTMiddleware memvalidasi token JWT.
Memperbarui konteks dengan user ID dan role jika valid.
*/
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("JWTMiddleware: Starting token validation for %s %s", r.Method, r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("JWTMiddleware: Authorization header missing")
			response.Unauthorized(w, "Authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)                  // Trim spasi ekstra
		log.Printf("JWTMiddleware: Extracted token: %q", tokenString) // Log token untuk debug
		if tokenString == "" || tokenString == authHeader {
			log.Printf("JWTMiddleware: Invalid token format")
			response.Unauthorized(w, "Invalid token format")
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			log.Printf("JWTMiddleware: Token validation failed: %v", err)
			response.Unauthorized(w, "Invalid token")
			return
		}

		log.Printf("JWTMiddleware: Token valid, userID: %s, role: %s", claims.Subject, claims.Role)
		log.Printf("JWTMiddleware: exp=%v (unix), now=%v", claims.ExpiresAt.Unix(), time.Now().Unix())
		// Atur konteks dengan userID dan role dari claims
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
