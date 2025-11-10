package service

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lalan-be/internal/config"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
)

/*
Merepresentasikan respons autentikasi.
Digunakan untuk mengembalikan ID, token akses, refresh, tipe, dan kedaluwarsa setelah registrasi/login sukses.
*/
type AuthResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
Mendefinisikan operasi service autentikasi.
Menyediakan method untuk registrasi, login, dan ambil profil dengan hasil sukses atau error.
*/
type AuthService interface {
	Register(input *model.HosterModel) (*AuthResponse, error)
	Login(email, password string) (*AuthResponse, error)
	GetHosterProfile(userID string) (*model.HosterModel, error)
}

/*
Implementasi service autentikasi dengan repository.
*/
type authService struct {
	repo repository.AuthRepository
}

/*
Membuat service autentikasi.
Mengembalikan instance AuthService yang siap digunakan.
*/
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

/*
Menghasilkan token JWT.
Mengembalikan AuthResponse dengan token akses dan refresh atau error jika gagal.
*/
func (s *authService) generateToken(userID string) (*AuthResponse, error) {
	// Access Token (1 jam)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	// Refresh Token (simpan di Redis nanti)
	refreshToken := uuid.New().String()

	return &AuthResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    16000,
	}, nil
}

/*
Mendaftarkan hoster baru.
Mengembalikan AuthResponse dengan token atau error jika registrasi gagal.
*/
func (s *authService) Register(input *model.HosterModel) (*AuthResponse, error) {
	// Generate UUID untuk ID
	input.ID = uuid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(message.MsgFailedToHashPassword)
	}
	input.PasswordHash = string(hashedPassword)

	// Simpan ke database
	if err := s.repo.CreateHoster(input); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "hosters_email_key") {
			return nil, errors.New(message.MsgHosterEmailExists)
		}
	}
	// Generate token (auto login)
	return s.generateToken(input.ID)
}

/*
Memproses login hoster.
Mengembalikan AuthResponse dengan token atau error jika kredensial salah.
*/
func (s *authService) Login(email, password string) (*AuthResponse, error) {
	hoster, err := s.repo.FindByEmailForLogin(email)
	if err != nil || hoster == nil {
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}

	return s.generateToken(hoster.ID)
}

/*
Mengambil profil hoster.
Mengembalikan data hoster atau error jika tidak ditemukan.
*/
func (s *authService) GetHosterProfile(userID string) (*model.HosterModel, error) {
	log.Printf("get prodile userid: %s", userID)
	hoster, err := s.repo.GetHosterByID(userID)
	if err != nil {
		log.Printf("Error getting profile: %v", err)
		return nil, errors.New(message.MsgInternalServerError)
	}
	if hoster == nil {
		return nil, errors.New(message.MsgHosterNotFound)
	}
	return hoster, nil
}
