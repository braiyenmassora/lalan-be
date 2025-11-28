package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"lalan-be/internal/config"
	"lalan-be/internal/message"
	"lalan-be/internal/response"

	"github.com/golang-jwt/jwt/v5"
)

/*
contextKey adalah tipe khusus untuk key pada context agar type-safe dan menghindari collision.
*/
type contextKey string

/*
UserIDKey dan UserRoleKey adalah key untuk menyimpan data user dari JWT ke dalam request context.
*/
const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

/*
Claims adalah custom JWT claims yang menyimpan role user selain standard claims.
*/
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/*
CORSMiddleware mengatur header CORS sesuai konfigurasi environment.

Alurures:
1. Baca ALLOWED_ORIGIN dari config
2. Log warning jika di production tapi masih pakai wildcard
3. Set header CORS
4. Tangani preflight request (OPTIONS)

Output:
- Lanjut ke handler berikutnya jika bukan OPTIONS
- 200 OK langsung untuk request OPTIONS
*/
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := config.GetEnv("APP_ENV", "dev")
		var allowedCfg string

		if env == "dev" {
			allowedCfg = config.GetEnv("ALLOWED_ORIGIN_DEV", "*")
		} else {
			allowedCfg = config.GetEnv("ALLOWED_ORIGIN_STAGING", "https://lalan-fe.vercel.app")
		}

		reqOrigin := r.Header.Get("Origin")
		allowed := "*"

		if env == "production" {
			// hanya izinkan origin yang match exact
			if reqOrigin == allowedCfg {
				allowed = reqOrigin
			} else {
				allowed = "" // tolak
			}
		} else {
			// dev/staging → mirror origin agar cocok dgn credentials
			if reqOrigin != "" {
				allowed = reqOrigin
			}
		}

		if allowed != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
GetUserID mengambil user ID yang sudah disimpan di context oleh JWTMiddleware.

Output sukses:
- string userID jika ada
- string kosong jika tidak ada / belum divalidasi
*/
func GetUserID(r *http.Request) string {
	if val := r.Context().Value(UserIDKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
GetUserRole mengambil role user dari context.

Output sukses:
- string role jika ada
- string kosong jika tidak ada / belum divalidasi
*/
func GetUserRole(r *http.Request) string {
	if val := r.Context().Value(UserRoleKey); val != nil {
		return val.(string)
	}
	return ""
}

/*
JWTMiddleware memvalidasi token JWT dari header Authorization dan mengisi context dengan user data.

Alur kerja:
1. Cek header Authorization (format: Bearer <token>)
2. Parse dan validasi token dengan secret key
3. Validasi signing method (hanya HS256)
4. Simpan user_id (subject) dan role ke context

Output sukses:
- Lanjut ke handler berikutnya dengan context yang sudah diisi
Output error:
- 401 Unauthorized → header kosong / format salah / token invalid / expired / algorithm salah
*/
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := config.GetEnv("APP_ENV", "dev")
		isDev := env != "production"

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if isDev {
				log.Printf("JWTMiddleware: missing Authorization header - %s %s", r.Method, r.URL.Path)
			}
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenStr == "" || tokenStr == authHeader {
			if isDev {
				log.Printf("JWTMiddleware: invalid Authorization format")
			}
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return config.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			if isDev {
				log.Printf("JWTMiddleware: token invalid - err: %v", err)
			}
			response.Unauthorized(w, message.Unauthorized)
			return
		}

		// Simpan ke context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		r = r.WithContext(ctx)

		if isDev {
			log.Printf("JWTMiddleware: authenticated user_id=%s role=%s", claims.Subject, claims.Role)
		}

		next.ServeHTTP(w, r)
	})
}

/*
Admin adalah middleware yang memastikan user memiliki role "admin".

Output sukses:
- Lanjut ke handler berikutnya
Output error:
- 403 Forbidden → role bukan admin
*/
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Admin middleware: checking access - user_id=%s role=%s", userID, role)

		if role != "admin" {
			response.Forbidden(w, message.AdminAccessRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*
Hoster adalah middleware yang memastikan user memiliki role "hoster".

Output sukses:
- Lanjut ke handler berikutnya
Output error:
- 403 Forbidden → role bukan hoster
*/
func Hoster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Hoster middleware: checking access - user_id=%s role=%s", userID, role)

		if role != "hoster" {
			response.Forbidden(w, message.HosterAccessRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*
Customer adalah middleware yang memastikan user memiliki role "customer".

Output sukses:
- Lanjut ke handler berikutnya
Output error:
- 403 Forbidden → role bukan customer
*/
func Customer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r)
		userID := GetUserID(r)
		log.Printf("Customer middleware: checking access - user_id=%s role=%s", userID, role)

		if role != "customer" {
			response.Forbidden(w, message.CustomerAccessRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*
GenerateToken membuat JWT baru dengan masa berlaku 1 jam.

Output sukses:
- (string token, nil)
Output error:
- ("", error) → gagal signing token
*/
func GenerateToken(userID, role string) (string, error) {
	expiration := time.Now().Add(1 * time.Hour)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTSecret())
}

/*
RefreshToken memvalidasi token lama dan mengeluarkan token baru dengan data yang sama.

Output sukses:
- (string newToken, nil)
Output error:
- ("", error) → token lama invalid / expired / malformed
*/
func RefreshToken(oldTokenString string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(oldTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return "", err
	}

	return GenerateToken(claims.Subject, claims.Role)
}
