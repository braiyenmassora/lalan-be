package service

import (
	"errors"
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
Struct untuk respons autentikasi.
*/
type AuthResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
Interface untuk operasi autentikasi.
*/
type AuthService interface {
	Register(input *model.HosterModel) (*AuthResponse, error)
	Login(email, password string) (*AuthResponse, error)
}

/*
Struct implementasi AuthService.
*/
type authService struct {
	repo repository.AuthRepository
}

/*
Membuat instance AuthService dengan repository.
*/
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

/*
Membuat token JWT untuk user.
Mengembalikan AuthResponse atau error.
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
Mengembalikan AuthResponse atau error.
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
Memvalidasi login hoster.
Mengembalikan AuthResponse atau error.
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
