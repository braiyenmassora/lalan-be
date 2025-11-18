package customer

import (
	"context"
	"errors"
	"lalan-be/internal/config"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/pkg/message"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
hosterService menyediakan logika bisnis untuk hoster.
Menggunakan repository untuk akses data.
*/
type customerService struct {
	repo CustomerRepository
}

/*
Methods untuk hosterService menangani operasi bisnis hoster, item, dan terms.
Dipanggil oleh handler untuk validasi dan logika.
*/
func (s *customerService) generateTokenCustomer(userID string) (*CustomerResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "customer",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return &CustomerResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *customerService) LoginCustomer(email, password string) (*CustomerResponse, error) {
	customer, err := s.repo.FindByEmailCustomerForLogin(email)
	if err != nil || customer == nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	return s.generateTokenCustomer(customer.ID)
}
func (s *customerService) CreateCustomer(customer *model.CustomerModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(customer.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	customer.PasswordHash = string(hash)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	err = s.repo.CreateCustomer(customer)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgCustomerEmailExists)
		}
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *customerService) GetDetailCustomer(ctx context.Context) (*model.CustomerModel, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	customer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if customer == nil {
		return nil, errors.New(message.MsgCustomerNotFound)
	}

	return customer, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, updateData *model.CustomerModel) error {
	// Ambil ID customer dari context (dari JWT token)
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	// Validasi: Pastikan hanya field yang diizinkan yang diubah
	// Ambil data customer yang ada untuk memastikan hanya field tertentu yang diupdate
	existingCustomer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if existingCustomer == nil {
		return errors.New(message.MsgCustomerNotFound)
	}

	// Update hanya field yang diizinkan: full_name, phone_number, profile_photo, address
	// Field lain (seperti email, password_hash) tetap dari existing data
	existingCustomer.FullName = updateData.FullName
	existingCustomer.PhoneNumber = updateData.PhoneNumber
	existingCustomer.ProfilePhoto = updateData.ProfilePhoto
	existingCustomer.Address = updateData.Address
	existingCustomer.UpdatedAt = time.Now() // Set waktu update sekarang

	// Panggil repository untuk update
	err = s.repo.UpdateCustomer(existingCustomer)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *customerService) DeleteCustomer(ctx context.Context) error {
	// Ambil ID customer dari context (dari JWT token)
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	// Opsional: Validasi apakah customer ada sebelum delete
	existingCustomer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if existingCustomer == nil {
		return errors.New(message.MsgCustomerNotFound)
	}

	// Panggil repository untuk delete
	err = s.repo.DeleteCustomer(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *customerService) UploadIdentity(ctx context.Context, ktpURL string) error {
	// Ambil ID customer dari context (dari JWT token)
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	// Opsional: Cek apakah customer sudah punya identity (untuk mencegah duplikat)
	// Jika perlu, tambahkan method di repo untuk cek berdasarkan user_id

	// Buat instance IdentityModel
	identity := &model.IdentityModel{
		UserID:         id,
		KTPURL:         ktpURL,
		Verified:       false,
		Status:         "pending",
		RejectedReason: "",
		VerifiedAt:     nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Panggil repository untuk insert
	err := s.repo.CreateIdentity(identity)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *customerService) CheckIdentityExists(ctx context.Context) error {
	// Ambil ID customer dari context
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	// Cek apakah identity sudah ada
	exists, err := s.repo.CheckIdentityExists(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if exists {
		return errors.New("Identity already exists")
	}
	return nil
}

func (s *customerService) GetIdentityStatus(ctx context.Context) (*model.IdentityModel, error) {
	// Ambil ID customer dari context
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	// Panggil repository
	identity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if identity == nil {
		return nil, errors.New("Identity not found")
	}

	return identity, nil
}

type CustomerResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type CustomerService interface {
	LoginCustomer(email, password string) (*CustomerResponse, error)
	CreateCustomer(customer *model.CustomerModel) error
	GetDetailCustomer(ctx context.Context) (*model.CustomerModel, error)
	UpdateCustomer(ctx context.Context, updateData *model.CustomerModel) error
	DeleteCustomer(ctx context.Context) error
	UploadIdentity(ctx context.Context, ktpURL string) error
	CheckIdentityExists(ctx context.Context) error
	GetIdentityStatus(ctx context.Context) (*model.IdentityModel, error) // Tambahkan ini
}

/*
NewCustomerService membuat instance CustomerService.
Menginisialisasi service dengan repository.
*/
func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}
