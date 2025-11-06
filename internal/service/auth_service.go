package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(input *model.HosterModel) error
	Login(email, password string) (*pkg.Response, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Register membuat akun hoster baru
func (s *authService) Register(input *model.HosterModel) error {
	// cek email sudah terdaftar atau belum
	existing, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(pkg.MsgHosterEmailExists)
	}

	// generate ID unik
	input.ID = uuid.New().String()

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	input.PasswordHash = string(hashed)

	// simpan ke database
	return s.repo.CreateHoster(input)
}

// Login hoster (cek email & password)
func (s *authService) Login(email, password string) (*pkg.Response, error) {
	// 1. Cari hoster berdasarkan email
	hoster, err := s.repo.FindByEmailForLogin(email)
	if err != nil {
		log.Printf("Login error: %v", err)
		return nil, err
	}
	if hoster == nil {
		return nil, errors.New(pkg.MsgHosterInvalidCredentials)
	}

	// 2. Bandingkan password (bcrypt)
	log.Printf("Login attempt: email=%s", email)
	if err := bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		return nil, errors.New(pkg.MsgHosterInvalidCredentials)
	}
	log.Println("Password match! Login success.")

	// 3. Clean hoster (hapus password)
	cleanHoster := *hoster
	cleanHoster.PasswordHash = ""

	// 4. Return response
	return &pkg.Response{
		Code:    http.StatusOK,
		Data:    map[string]interface{}{"hoster": cleanHoster},
		Message: pkg.MsgHosterLoginSuccess,
		Status:  "success",
	}, nil
}
