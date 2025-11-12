package service

import (
	"context"
	"errors"
	"lalan-be/internal/config"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CustomerResponse merepresentasikan respons autentikasi customer.
// Digunakan untuk mengembalikan ID, token akses, refresh, tipe, dan kedaluwarsa setelah registrasi/login sukses.
type CustomerResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// CustomerClaim custom claims dengan role untuk customer.
type CustomerClaim struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// customerService mengimplementasikan service autentikasi dengan repository.
type customerService struct {
	repo repository.CustomerRepository
}

// generateTokenCustomer menghasilkan token JWT untuk customer.
// Mengembalikan CustomerResponse dengan token akses dan refresh atau error jika gagal.
func (s *customerService) generateTokenCustomer(userID string) (*CustomerResponse, error) {
	// Access Token (1 jam)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := CustomerClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "customer",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	// Refresh Token (simpan di Redis nanti)
	refreshToken := uuid.New().String()

	return &CustomerResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    16000,
	}, nil
}

// RegisterCustomer mendaftarkan customer baru.
// Mengembalikan CustomerResponse dengan token atau error jika registrasi gagal.
func (s *customerService) RegisterCustomer(input *model.CustomerModel) (*CustomerResponse, error) {
	// Generate UUID untuk ID
	input.ID = uuid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(message.MsgFailedToHashPassword)
	}
	input.PasswordHash = string(hashedPassword)

	// Simpan ke database
	if err := s.repo.CreateCustomer(input); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "customers_email_key") {
			return nil, errors.New(message.MsgCustomerEmailExists)
		}
	}
	// Generate token (auto login)
	return s.generateTokenCustomer(input.ID)
}

// LoginCustomer memproses login customer.
// Mengembalikan CustomerResponse dengan token atau error jika kredensial salah.
func (s *customerService) LoginCustomer(email, password string) (*CustomerResponse, error) {
	customer, err := s.repo.FindByEmailCustomerForLogin(email)
	if err != nil || customer == nil {
		return nil, errors.New(message.MsgCustomerInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New(message.MsgCustomerInvalidCredentials)
	}

	return s.generateTokenCustomer(customer.ID)
}

// GetCustomerProfile mengambil profil customer berdasarkan user yang sedang login (dari token).
// Mengembalikan data customer atau error jika tidak ditemukan.
func (s *customerService) GetCustomerProfile(ctx context.Context) (*model.CustomerModel, error) {
	// Ekstrak userID dari context (misalnya dari JWT claims yang disimpan di middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errors.New(message.MsgUnauthorized)
	}

	log.Printf("get profile userid: %s", userID)
	customer, err := s.repo.GetCustomerByID(userID)
	if err != nil {
		log.Printf("Error getting profile: %v", err)
		return nil, errors.New(message.MsgInternalServerError)
	}
	if customer == nil {
		return nil, errors.New(message.MsgCustomerNotFound)
	}
	return customer, nil
}

// CustomerService mendefinisikan operasi service autentikasi customer.
// Menyediakan method untuk registrasi, login, dan ambil profil dengan hasil sukses atau error.
type CustomerService interface {
	RegisterCustomer(input *model.CustomerModel) (*CustomerResponse, error)
	LoginCustomer(email, password string) (*CustomerResponse, error)
	GetCustomerProfile(ctx context.Context) (*model.CustomerModel, error) // Ubah parameter ke context
}

// NewCustomerService membuat service autentikasi.
// Mengembalikan instance CustomerService yang siap digunakan.
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}
