package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService adalah interface untuk layanan autentikasi
type AuthService interface {
	Register(input *model.HosterModel) error                  // Mendaftarkan hoster baru
	Login(email, password string) (*response.Response, error) // Masuk dengan email dan password
}

// authService adalah struct yang mengimplementasikan AuthService
type authService struct {
	repo repository.AuthRepository // Repository untuk operasi autentikasi
}

// NewAuthService membuat instance baru dari authService
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Register mendaftarkan hoster baru
func (s *authService) Register(input *model.HosterModel) error {
	// Memeriksa apakah email sudah terdaftar
	existing, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(message.MsgHosterEmailExists)
	}

	// Menghasilkan ID unik
	input.ID = uuid.New().String()

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	input.PasswordHash = string(hashed)

	// Menyimpan ke database
	return s.repo.CreateHoster(input)
}

// Login masuk dengan email dan password
func (s *authService) Login(email, password string) (*response.Response, error) {
	// Mencari hoster berdasarkan email
	hoster, err := s.repo.FindByEmailForLogin(email)
	if err != nil {
		log.Printf("Login error: %v", err)
		return nil, err
	}
	if hoster == nil {
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}

	// Membandingkan password
	log.Printf("Login attempt: email=%s", email)
	if err := bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}
	log.Println("Password match! Login success.")

	// Membersihkan data hoster
	cleanHoster := *hoster
	cleanHoster.PasswordHash = ""

	// Mengembalikan response
	return &response.Response{
		Code:    http.StatusOK,
		Data:    map[string]interface{}{"hoster": cleanHoster},
		Message: message.MsgHosterCreatedSuccess,
		Status:  "success",
	}, nil
}
