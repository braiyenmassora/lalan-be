package admin

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lalan-be/internal/config"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/pkg/message"
)

/*
AdminResponse berisi data token dan informasi admin.
Digunakan untuk respons autentikasi admin.
*/
type AdminResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
adminService menyediakan logika bisnis untuk operasi admin.
Menggunakan repository untuk akses data.
*/
type adminService struct {
	repo AdminRepository
}

/*
Methods adminService menangani autentikasi dan manajemen admin serta kategori.
Menggunakan bcrypt untuk hashing dan JWT untuk token.
*/
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
		return nil, errors.New(message.MsgInternalServerError)
	}

	return &AdminResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *adminService) LoginAdmin(email, password string) (*AdminResponse, error) {
	admin, err := s.repo.FindByEmailAdminForLogin(email)
	if err != nil || admin == nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	return s.generateTokenAdmin(admin.ID)
}

func (s *adminService) CreateAdmin(admin *model.AdminModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	admin.PasswordHash = string(hash)
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	err = s.repo.CreateAdmin(admin)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgBadRequest)
		}
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *adminService) CreateCategory(category *model.CategoryModel) error {
	existing, err := s.repo.FindCategoryByName(category.Name)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if existing != nil {
		return errors.New(message.MsgBadRequest)
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	err = s.repo.CreateCategory(category)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *adminService) UpdateCategory(category *model.CategoryModel) error {
	category.UpdatedAt = time.Now()
	err := s.repo.UpdateCategory(category)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	return nil
}

func (s *adminService) DeleteCategory(id string) error {
	err := s.repo.DeleteCategory(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	return nil
}

func (s *adminService) GetAllCategory() ([]*model.CategoryModel, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return categories, nil
}

/*
AdminService mendefinisikan kontrak untuk logika bisnis admin.
Wajib diimplementasikan oleh semua penyedia layanan admin.
*/
type AdminService interface {
	CreateAdmin(*model.AdminModel) error
	LoginAdmin(email, password string) (*AdminResponse, error)
	CreateCategory(*model.CategoryModel) error
	UpdateCategory(*model.CategoryModel) error
	DeleteCategory(id string) error
	GetAllCategory() ([]*model.CategoryModel, error)
}

/*
NewAdminService membuat instance baru AdminService.
Mengembalikan interface AdminService dengan repository yang diberikan.
*/
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}
