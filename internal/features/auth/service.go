package auth

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lalan-be/internal/config"
	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
)

// AuthResponse sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
authService menangani logika bisnis untuk fitur autentikasi.
Struct ini menjadi penghubung antara handler dan repository.
*/
type authService struct {
	repo *authRepository
}

/*
NewAuthService membuat instance service baru.

Output:
- Pointer ke authService yang siap digunakan.
*/
func NewAuthService(repo *authRepository) *authService {
	return &authService{repo: repo}
}

/*
generateToken membuat JWT token baru untuk user yang berhasil login.

Fungsi ini membuat claims berisi UserID dan Role, lalu menandatanganinya
menggunakan secret key dari konfigurasi.

Output:
- Pointer ke AuthResponse berisi token dan metadata.
- error jika signing token gagal.
*/
func (s *authService) generateToken(userID, role string) (*dto.AuthResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return &dto.AuthResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Role:         role,
	}, nil
}

// CreateCustomerResponse sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
generateOTP membuat kode OTP numerik 6 digit secara acak.

Output:
- String berisi 6 digit angka (contoh: "123456").
*/
func (s *authService) generateOTP() string {
	const otpChars = "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += string(otpChars[num.Int64()])
	}
	return otp
}

/*
RegisterHoster mendaftarkan hoster baru.

Langkah-langkah:
1. Hash password.
2. Simpan data hoster ke database.

Output:
- error jika terjadi kesalahan sistem.
- nil jika berhasil.
*/
func (s *authService) RegisterHoster(h *domain.Hoster) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(h.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}
	h.PasswordHash = string(hash)
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	if err := s.repo.CreateHoster(h); err != nil {
		if err.Error() == "email already exists" {
			return errors.New(message.EmailAlreadyExists)
		}
		return errors.New(message.InternalError)
	}
	return nil
}

/*
RegisterAdmin mendaftarkan admin baru.

Langkah-langkah:
1. Hash password.
2. Simpan data admin ke database.

Output:
- error jika email duplikat atau kesalahan sistem.
- nil jika berhasil.
*/
func (s *authService) RegisterAdmin(a *domain.Admin) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}
	a.PasswordHash = string(hash)
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	if err := s.repo.CreateAdmin(a); err != nil {
		if err.Error() == "duplicate" {
			return errors.New(message.EmailAlreadyExists)
		}
		return errors.New(message.InternalError)
	}
	return nil
}

/*
RegisterCustomer mendaftarkan customer baru.

Langkah-langkah:
1. Hash password.
2. Generate OTP dan set expiry.
3. Simpan data customer ke database.

Output:
- error jika email duplikat atau kesalahan sistem.
- nil jika berhasil.
*/
func (s *authService) RegisterCustomer(c *domain.Customer) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}
	c.PasswordHash = string(hash)
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	// Generate OTP untuk verifikasi email
	otp := s.generateOTP()
	c.VerificationToken = otp
	c.VerificationExpiresAt = time.Now().Add(5 * time.Minute)

	if err := s.repo.CreateCustomer(c); err != nil {
		if err.Error() == "email already exists" {
			return errors.New(message.EmailAlreadyExists)
		}
		return errors.New(message.InternalError)
	}

	// TODO: Kirim email OTP untuk verifikasi

	return nil
}

/*
SendOTP memverifikasi kode OTP yang dikirimkan user.

Output:
- error jika OTP salah/kadaluarsa.
- nil jika verifikasi berhasil.
*/
func (s *authService) SendOTP(email string, otp string) error {
	if err := s.repo.SendOTP(email, otp); err != nil {
		if err.Error() == message.OTPInvalid {
			return errors.New(message.OTPInvalid)
		}
		return errors.New(message.InternalError)
	}
	return nil
}

/*
ResendOTP mengirim ulang kode OTP baru.

Langkah-langkah:
1. Generate OTP baru.
2. Update database dengan OTP baru dan expiry time baru (hanya jika belum verified).
3. (TODO) Kirim email.

Output:
- String OTP baru (untuk keperluan dev/testing).
- error jika customer tidak ditemukan atau sudah verified.
*/
func (s *authService) ResendOTP(email string) (string, error) {
	newOTP := s.generateOTP()
	exp := time.Now().Add(5 * time.Minute)

	if err := s.repo.ResendOTP(email, newOTP, exp); err != nil {
		if err.Error() == message.CustomerNotFound {
			return "", errors.New(message.CustomerNotFound)
		}
		if err.Error() == message.OTPAlreadyVerified {
			return "", errors.New(message.OTPAlreadyVerified)
		}
		return "", errors.New(message.InternalError)
	}

	// TODO: Kirim email OTP baru

	return newOTP, nil
}

/*
Login memproses login untuk semua role (admin, hoster, customer).

Langkah-langkah:
1. Cari user berdasarkan email di semua tabel.
2. Jika user ditemukan, cek password hash.
3. Khusus customer, cek apakah email sudah diverifikasi.
4. Jika valid, generate JWT token.

Output:
- Pointer ke AuthResponse berisi token.
- error jika login gagal (user tidak ditemukan, password salah, email belum verifikasi).
*/
func (s *authService) Login(email, password string) (*dto.AuthResponse, error) {
	// 1. Cari user
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if user == nil {
		return nil, errors.New(message.LoginFailed)
	}

	// 2. Cek verifikasi email (khusus customer)
	if user.Role == "customer" && !user.EmailVerified {
		return nil, errors.New(message.EmailNotVerified)
	}

	// 3. Verifikasi password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.LoginFailed)
	}

	// 4. Generate token
	return s.generateToken(user.ID, user.Role)
}

/*
ForgotPassword generates reset token dan kirim via email.

Output:
- String reset token (untuk dev/testing).
- error jika user tidak ditemukan.
*/
func (s *authService) ForgotPassword(email, role string) (string, error) {
	// Validasi role
	if role != "customer" && role != "hoster" {
		return "", errors.New("invalid role")
	}

	// Generate reset token
	resetToken := s.generateOTP()
	exp := time.Now().Add(15 * time.Minute)

	if err := s.repo.RequestPasswordReset(email, role, resetToken, exp); err != nil {
		if err.Error() == message.CustomerNotFound || err.Error() == message.HosterNotFound {
			return "", err
		}
		return "", errors.New(message.InternalError)
	}

	// TODO: Kirim email reset token

	return resetToken, nil
}

/*
ResetPassword verifikasi reset token dan update password baru.

Output:
- error jika token invalid/expired atau update gagal.
- nil jika berhasil.
*/
func (s *authService) ResetPassword(email, role, token, newPassword string) error {
	// Validasi role
	if role != "customer" && role != "hoster" {
		return errors.New("invalid role")
	}

	// Hash password baru
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}

	// Update password
	if err := s.repo.ResetPassword(email, role, token, string(hash)); err != nil {
		if err.Error() == message.ResetTokenInvalid {
			return errors.New(message.ResetTokenInvalid)
		}
		return errors.New(message.InternalError)
	}

	return nil
}
