package admin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lalan-be/internal/config"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
)

// AdminResponse berisi data token dan informasi admin untuk respons autentikasi
type AdminResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// adminService menyediakan logika bisnis untuk operasi admin
type adminService struct {
	repo AdminRepository
}

// generateTokenAdmin menghasilkan token JWT untuk admin dengan klaim tertentu
func (s *adminService) generateTokenAdmin(userID string) (*AdminResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return &AdminResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

// LoginAdmin melakukan autentikasi admin dan menghasilkan token jika berhasil
func (s *adminService) LoginAdmin(email, password string) (*AdminResponse, error) {
	admin, err := s.repo.FindByEmailAdminForLogin(email)
	if err != nil || admin == nil {
		return nil, errors.New(message.LoginFailed)
	}

	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.LoginFailed)
	}

	return s.generateTokenAdmin(admin.ID)
}

// CreateAdmin membuat admin baru dengan hashing password dan penyimpanan data
func (s *adminService) CreateAdmin(admin *model.AdminModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}
	admin.PasswordHash = string(hash)
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	err = s.repo.CreateAdmin(admin)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(fmt.Sprintf(message.AlreadyExists, "admin"))
		}
		return errors.New(message.InternalError)
	}

	return nil
}

// CreateCategory membuat kategori baru setelah validasi nama unik
func (s *adminService) CreateCategory(category *model.CategoryModel) error {
	existing, err := s.repo.FindCategoryByName(category.Name)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existing != nil {
		return errors.New(fmt.Sprintf(message.AlreadyExists, "category"))
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	err = s.repo.CreateCategory(category)
	if err != nil {
		return errors.New(message.InternalError)
	}

	return nil
}

// UpdateCategory memperbarui data kategori berdasarkan ID
func (s *adminService) UpdateCategory(category *model.CategoryModel) error {
	category.UpdatedAt = time.Now()
	err := s.repo.UpdateCategory(category)
	if err != nil {
		return errors.New(message.InternalError)
	}
	return nil
}

// DeleteCategory menghapus kategori berdasarkan ID
func (s *adminService) DeleteCategory(id string) error {
	err := s.repo.DeleteCategory(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	return nil
}

// GetAllCategory mengambil semua data kategori
func (s *adminService) GetAllCategory() ([]*model.CategoryModel, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return categories, nil
}

/*
UpdateIdentityStatus
memperbarui status identitas berdasarkan user ID
*/
func (s *adminService) UpdateIdentityStatus(ctx context.Context, userID string, status string, rejectedReason string) error {
	// Ambil ID admin dari context
	adminID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		log.Printf("UpdateIdentityStatus: unauthorized access")
		return errors.New(message.Unauthorized)
	}

	log.Printf("UpdateIdentityStatus: admin %s updating identity for user %s to status %s", adminID, userID, status)

	// Validasi status
	if status != "approved" && status != "rejected" {
		log.Printf("UpdateIdentityStatus: invalid status %s", status)
		return errors.New(message.InvalidStatus)
	}

	// Cek apakah identity ada berdasarkan user ID
	identity, err := s.repo.GetIdentityByCustomerID(userID)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error getting identity for user %s: %v", userID, err)
		return errors.New(message.InternalError)
	}
	if identity == nil {
		log.Printf("UpdateIdentityStatus: identity not found for user %s", userID)
		return errors.New(fmt.Sprintf(message.NotFound, "identity"))
	}

	// Hitung verified dan verifiedAt
	var verified bool
	var verifiedAt *time.Time
	if status == "approved" {
		verified = true
		now := time.Now()
		verifiedAt = &now
	} else {
		verified = false
		verifiedAt = nil
	}

	// Update status menggunakan identity ID
	err = s.repo.UpdateIdentityStatus(identity.ID, status, rejectedReason, verified, verifiedAt)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error updating identity %s: %v", identity.ID, err)
		return errors.New(message.InternalError)
	}

	log.Printf("UpdateIdentityStatus: successfully updated identity for user %s to status %s", userID, status)
	return nil
}

/*
GetIdentityByCustomerID
mengambil data identitas berdasarkan user ID
*/
func (s *adminService) GetIdentityByCustomerID(userID string) (*model.IdentityModel, error) {
	identity, err := s.repo.GetIdentityByCustomerID(userID)
	if err != nil {
		log.Printf("GetIdentityByCustomerID: error getting identity for user %s: %v", userID, err)
		return nil, errors.New(message.InternalError)
	}
	if identity == nil {
		log.Printf("GetIdentityByCustomerID: identity not found for user %s", userID)
		return nil, errors.New(fmt.Sprintf(message.NotFound, "identity"))
	}
	return identity, nil
}

// AdminService mendefinisikan kontrak untuk logika bisnis admin
type AdminService interface {
	CreateAdmin(*model.AdminModel) error
	LoginAdmin(email, password string) (*AdminResponse, error)
	CreateCategory(*model.CategoryModel) error
	UpdateCategory(*model.CategoryModel) error
	DeleteCategory(id string) error
	GetAllCategory() ([]*model.CategoryModel, error)
	UpdateIdentityStatus(context.Context, string, string, string) error
	GetIdentityByCustomerID(string) (*model.IdentityModel, error)
}

// NewAdminService membuat instance baru AdminService dengan repository yang diberikan
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}
