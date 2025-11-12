package service

import (
	"context"
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
Merepresentasikan respons autentikasi hoster.
Digunakan untuk mengembalikan ID, token akses, refresh, tipe, dan kedaluwarsa setelah registrasi/login sukses.
*/
type HosterResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
Custom claims dengan role untuk hoster.
*/
type HosterClaim struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/*
Implementasi service autentikasi dengan repository.
*/
type hosterService struct {
	repo repository.AuthRepository
}

/*
Menghasilkan token JWT untuk hoster.
Mengembalikan HosterResponse dengan token akses dan refresh atau error jika gagal.
*/
func (s *hosterService) generateTokenHoster(userID string) (*HosterResponse, error) {
	// Access Token (1 jam)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := HosterClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "hoster",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	// Refresh Token (simpan di Redis nanti)
	refreshToken := uuid.New().String()

	return &HosterResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    16000,
	}, nil
}

/*
Mendaftarkan hoster baru.
Mengembalikan HosterResponse dengan token atau error jika registrasi gagal.
*/
func (s *hosterService) Register(input *model.HosterModel) (*HosterResponse, error) {
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
	return s.generateTokenHoster(input.ID)
}

/*
Memproses login hoster.
Mengembalikan HosterResponse dengan token atau error jika kredensial salah.
*/
func (s *hosterService) LoginHoster(email, password string) (*HosterResponse, error) {
	hoster, err := s.repo.FindByEmailForLogin(email)
	if err != nil || hoster == nil {
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New(message.MsgHosterInvalidCredentials)
	}

	return s.generateTokenHoster(hoster.ID)
}

/*
Mengambil profil hoster berdasarkan user yang sedang login (dari token).
Mengembalikan data hoster atau error jika tidak ditemukan.
*/
func (s *hosterService) GetHosterProfile(ctx context.Context) (*model.HosterModel, error) {
	// Ekstrak userID dari context (misalnya dari JWT claims yang disimpan di middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errors.New(message.MsgUnauthorized)
	}

	log.Printf("get profile userid: %s", userID)
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

/*
Mendefinisikan operasi service autentikasi hoster.
Menyediakan method untuk registrasi, login, dan ambil profil dengan hasil sukses atau error.
*/
type HosterService interface {
	Register(input *model.HosterModel) (*HosterResponse, error)
	LoginHoster(email, password string) (*HosterResponse, error)
	GetHosterProfile(ctx context.Context) (*model.HosterModel, error) // Ubah parameter ke context
}

/*
Membuat service autentikasi.
Mengembalikan instance HosterService yang siap digunakan.
*/
func NewHosterService(repo repository.AuthRepository) HosterService {
	return &hosterService{repo: repo}
}
