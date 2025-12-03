package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/response"
)

/*
AuthHandler menangani endpoint autentikasi terpusat.
Struct ini bertanggung jawab untuk memproses request login, register,
dan verifikasi OTP untuk semua role user (admin, hoster, customer).
*/
type AuthHandler struct {
	service *authService
}

/*
NewAuthHandler membuat instance handler baru.

Output:
- Pointer ke AuthHandler yang siap digunakan.
*/
func NewAuthHandler(s *authService) *AuthHandler {
	return &AuthHandler{service: s}
}

// LoginRequest sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
Login menangani proses autentikasi user.

Fungsi ini memvalidasi kredensial user (email & password),
memeriksa format email, dan jika valid, mengembalikan token akses.

Langkah-langkah:
1. Validasi method request (harus POST).
2. Decode & validasi JSON body.
3. Validasi format email menggunakan regex.
4. Panggil service login untuk verifikasi kredensial.
5. Set cookie auth_token (HttpOnly).
6. Kembalikan response sukses dengan data user & token.

Output:
- 200 OK: Login berhasil, return user data & token.
- 400 Bad Request: Input tidak valid atau email belum diverifikasi.
- 401 Unauthorized: Email atau password salah.
*/
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.Login: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Auth.Login: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Printf("Auth.Login: email or password empty")
		response.BadRequest(w, message.BadRequest)
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		log.Printf("Auth.Login: invalid email format: %s", req.Email)
		response.BadRequest(w, fmt.Sprintf(message.InvalidFormat, "email"))
		return
	}

	resp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("Auth.Login: login failed: %v", err)
		// Customer email not verified
		if err.Error() == message.EmailNotVerified {
			response.BadRequest(w, message.EmailNotVerified)
			return
		}
		response.Unauthorized(w, message.LoginFailed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    resp.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   resp.ExpiresIn,
	})

	userData := map[string]interface{}{
		"id":            resp.ID,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"token_type":    resp.TokenType,
		"expires_in":    resp.ExpiresIn,
		"role":          resp.Role,
	}

	response.OK(w, userData, message.Success)
}

// RegisterRequest sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
Register menangani pendaftaran user baru.

Fungsi ini mendukung pendaftaran untuk 3 role berbeda:
1. Customer: Butuh verifikasi email (OTP).
2. Hoster: Langsung aktif (bisa disesuaikan).
3. Admin: Langsung aktif.

Langkah-langkah:
1. Validasi method & JSON body.
2. Cek role yang diminta.
3. Validasi field wajib sesuai role.
4. Panggil service register yang sesuai.
5. Kembalikan response sukses.

Output:
- 200 OK: Registrasi berhasil (atau OTP terkirim untuk customer).
- 400 Bad Request: Input tidak valid atau email sudah terdaftar.
- 500 Internal Server Error: Kesalahan server database/sistem.
*/
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.Register: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.RegisterRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("Auth.Register: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	// Handle Customer Registration
	if req.Role == "customer" {
		if req.FullName == "" || req.Email == "" || req.Password == "" {
			response.BadRequest(w, message.BadRequest)
			return
		}

		cust := &domain.Customer{
			FullName:     req.FullName,
			Email:        req.Email,
			PasswordHash: req.Password,
			PhoneNumber:  req.PhoneNumber,
			Address:      req.Address,
			ProfilePhoto: req.ProfilePhoto,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := h.service.RegisterCustomer(cust); err != nil {
			log.Printf("Auth.Register: create customer failed: %v", err)
			if err.Error() == message.EmailAlreadyExists {
				response.BadRequest(w, message.EmailAlreadyExists)
				return
			}
			response.Error(w, http.StatusInternalServerError, message.InternalError)
			return
		}
		response.OK(w, nil, message.CustomerCreated)
		return
	}

	// Handle Hoster Registration
	if req.Role == "hoster" {
		if req.FullName == "" || req.Email == "" || req.Password == "" {
			response.BadRequest(w, message.BadRequest)
			return
		}
		hoster := &domain.Hoster{
			FullName:     req.FullName,
			StoreName:    req.StoreName,
			PhoneNumber:  req.PhoneNumber,
			Email:        req.Email,
			PasswordHash: req.Password,
			Address:      req.Address,
			ProfilePhoto: req.ProfilePhoto,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := h.service.RegisterHoster(hoster); err != nil {
			log.Printf("Auth.Register: create hoster failed: %v", err)
			if err.Error() == message.EmailAlreadyExists {
				response.BadRequest(w, message.EmailAlreadyExists)
				return
			}
			response.Error(w, http.StatusInternalServerError, message.InternalError)
			return
		}
		response.OK(w, hoster, message.HosterCreated)
		return
	}

	// Handle Admin Registration
	if req.Role == "admin" {
		if req.FullName == "" || req.Email == "" || req.Password == "" {
			response.BadRequest(w, message.BadRequest)
			return
		}
		admin := &domain.Admin{
			FullName:     req.FullName,
			Email:        req.Email,
			PasswordHash: req.Password,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := h.service.RegisterAdmin(admin); err != nil {
			log.Printf("Auth.Register: create admin failed: %v", err)
			if err.Error() == message.EmailAlreadyExists {
				response.BadRequest(w, message.EmailAlreadyExists)
				return
			}
			response.Error(w, http.StatusInternalServerError, message.InternalError)
			return
		}
		response.OK(w, admin, message.AdminCreated)
		return
	}

	response.BadRequest(w, message.BadRequest)
}

// VerifyEmailRequest sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
VerifyEmail memverifikasi kode OTP untuk aktivasi akun customer.

Fungsi ini menerima email dan kode OTP, lalu memvalidasinya ke database.
Jika valid, status akun customer akan diubah menjadi aktif (email_verified = true).

Output:
- 200 OK: Verifikasi berhasil.
- 400 Bad Request: OTP salah atau kadaluarsa.
- 500 Internal Server Error: Kesalahan server.
*/
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.VerifyEmail: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.VerifyEmailRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("Auth.VerifyEmail: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.Email == "" || req.OTP == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	if err := h.service.SendOTP(req.Email, req.OTP); err != nil {
		log.Printf("Auth.VerifyEmail: error: %v", err)
		if err.Error() == message.OTPInvalid {
			response.BadRequest(w, message.OTPInvalid)
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, nil, message.Success)
}

// ResendOTPRequest sekarang menggunakan DTO dari package dto
// Lihat: internal/dto/auth_dto.go

/*
ResendOTP mengirim ulang kode OTP ke email customer.

Digunakan jika customer tidak menerima kode OTP sebelumnya atau
kode sudah kadaluarsa. Akan generate kode baru dan kirim via email.

Output:
- 200 OK: OTP baru berhasil dikirim (mengembalikan kode OTP untuk dev purposes).
- 400 Bad Request: Email tidak valid.
- 500 Internal Server Error: Gagal generate atau kirim email.
*/
func (h *AuthHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.ResendOTP: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.ResendOTPRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("Auth.ResendOTP: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.Email == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	otp, err := h.service.ResendOTP(req.Email)
	if err != nil {
		log.Printf("Auth.ResendOTP: error: %v", err)
		if err.Error() == message.OTPAlreadyVerified {
			response.BadRequest(w, message.OTPAlreadyVerified)
			return
		}
		if err.Error() == message.CustomerNotFound {
			response.BadRequest(w, message.CustomerNotFound)
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	// Note: Mengembalikan OTP di response hanya untuk development/testing
	response.OK(w, map[string]string{"otp": otp}, message.Success)
}

/*
ForgotPassword menangani request lupa password untuk customer dan hoster.

Fungsi ini generate reset token dan kirim ke email user.

Output:
- 200 OK: Reset token berhasil dikirim (return token untuk dev).
- 400 Bad Request: Email tidak valid atau role tidak valid.
- 500 Internal Server Error: Kesalahan server.
*/
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.ForgotPassword: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.ForgotPasswordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("Auth.ForgotPassword: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.Email == "" || req.Role == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	token, err := h.service.ForgotPassword(req.Email, req.Role)
	if err != nil {
		log.Printf("Auth.ForgotPassword: error: %v", err)
		if err.Error() == message.CustomerNotFound {
			response.BadRequest(w, message.CustomerNotFound)
			return
		}
		if err.Error() == message.HosterNotFound {
			response.BadRequest(w, message.HosterNotFound)
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	// Note: Mengembalikan token di response hanya untuk development/testing
	response.OK(w, map[string]string{"reset_token": token}, message.ResetTokenSent)
}

/*
ResetPassword menangani reset password dengan token.

Fungsi ini memverifikasi token dan update password baru.

Output:
- 200 OK: Password berhasil direset.
- 400 Bad Request: Token invalid/expired atau input tidak valid.
- 500 Internal Server Error: Kesalahan server.
*/
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("Auth.ResetPassword: received request")
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	var req dto.ResetPasswordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Printf("Auth.ResetPassword: invalid JSON: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}
	if req.Email == "" || req.Role == "" || req.Token == "" || req.NewPassword == "" {
		response.BadRequest(w, message.BadRequest)
		return
	}

	if err := h.service.ResetPassword(req.Email, req.Role, req.Token, req.NewPassword); err != nil {
		log.Printf("Auth.ResetPassword: error: %v", err)
		if err.Error() == message.ResetTokenInvalid {
			response.BadRequest(w, message.ResetTokenInvalid)
			return
		}
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}
	response.OK(w, nil, message.PasswordResetSuccess)
}
