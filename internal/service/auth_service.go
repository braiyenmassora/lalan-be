package service

import (
	"errors"
	"lalan-be/internal/config"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Paket service untuk handle logika bisnis autentikasi.

// Struktur respons autentikasi untuk API.
type AuthResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Interface untuk layanan autentikasi.
type AuthService interface {
	Register(input *model.HosterModel) (*AuthResponse, error)
	Login(email, password string) (*AuthResponse, error)
}

// Struktur implementasi layanan autentikasi.
type authService struct {
	repo repository.AuthRepository
}

// Buat instance baru layanan autentikasi dengan dependency injection.
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Generate token JWT akses dan refresh untuk user.
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
		ExpiresIn:    3600,
	}, nil
}

// Daftarkan hoster baru dan kembalikan token autentikasi.
func (s *authService) Register(input *model.HosterModel) (*AuthResponse, error) {
	// Generate UUID untuk ID
	input.ID = uuid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	input.PasswordHash = string(hashedPassword)

	// Simpan ke database
	if err := s.repo.CreateHoster(input); err != nil {
		return nil, err
	}

	// Generate token (auto login)
	return s.generateToken(input.ID)
}

// Validasi kredensial dan kembalikan token akses jika berhasil.
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
