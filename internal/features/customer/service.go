package customer

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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

/*
customerService
mengelola logika bisnis untuk customer menggunakan repository
*/
type customerService struct {
	repo CustomerRepository
}

/*
generateTokenCustomer
menghasilkan JWT token untuk customer
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
		return nil, errors.New(message.InternalError)
	}

	return &CustomerResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

/*
LoginCustomer
memvalidasi login customer dan menghasilkan token
*/
func (s *customerService) LoginCustomer(email, password string) (*CustomerResponse, error) {
	customer, err := s.repo.FindByEmailCustomerForLogin(email)
	if err != nil || customer == nil {
		return nil, errors.New(message.LoginFailed)
	}

	if !customer.EmailVerified {
		return nil, errors.New("Email belum diverifikasi. Silakan verifikasi email terlebih dahulu.")
	}

	if bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.LoginFailed)
	}

	return s.generateTokenCustomer(customer.ID)
}

/*
CreateCustomer
membuat customer baru dengan hash password
*/
func (s *customerService) CreateCustomer(customer *model.CustomerModel) (*CreateCustomerResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("CreateCustomer: error hashing password: %v", err)
		return nil, errors.New(message.InternalError)
	}
	customer.PasswordHash = string(hashedPassword)

	// Generate OTP
	otp := s.generateOTP()

	// Set email verification fields
	customer.EmailVerified = false
	customer.VerificationToken = otp
	customer.VerificationExpiresAt = &time.Time{}
	*customer.VerificationExpiresAt = time.Now().Add(5 * time.Minute)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	err = s.repo.CreateCustomer(customer)
	if err != nil {
		if err.Error() == "email already exists" {
			return nil, errors.New(message.EmailAlreadyExists)
		}
		log.Printf("CreateCustomer: error creating customer: %v", err)
		return nil, errors.New(message.InternalError)
	}

	// TODO: Send OTP email (integrate with email service)

	log.Printf("CreateCustomer: customer created with email %s", customer.Email)
	return &CreateCustomerResponse{
		Customer: customer,
		OTP:      otp,
	}, nil
}

/*
generateOTP
menghasilkan OTP 6 digit
*/
func (s *customerService) generateOTP() string {
	const otpChars = "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += string(otpChars[num.Int64()])
	}
	return otp
}

/*
SendOTP
memverifikasi email customer dengan OTP
*/
func (s *customerService) SendOTP(email string, otp string) error {
	err := s.repo.SendOTP(email, otp)
	if err != nil {
		if err.Error() == "invalid or expired OTP" {
			return errors.New(message.OTPInvalid)
		}
		return errors.New(message.InternalError)
	}
	log.Printf("SendOTP: email %s verified successfully", email)
	return nil
}

/*
ResendOTP
mengirim ulang OTP dengan token baru dan waktu kadaluarsa
*/
func (s *customerService) ResendOTP(email string) (*ResendOTPResponse, error) {
	// Check if customer exists and if email is already verified
	customer, err := s.repo.FindByEmailCustomerForLogin(email)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if customer == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "customer"))
	}
	if customer.EmailVerified {
		return nil, errors.New(message.OTPAlreadyVerified)
	}

	// Generate new OTP
	newOTP := s.generateOTP()

	// Set expiration time (e.g., 5 minutes from now)
	expireTime := time.Now().Add(5 * time.Minute)

	// Update customer with new OTP via repository
	err = s.repo.ResendOTP(email, newOTP, expireTime)
	if err != nil {
		if err.Error() == "customer not found" {
			return nil, errors.New(fmt.Sprintf(message.NotFound, "customer"))
		}
		return nil, errors.New(message.InternalError)
	}

	// TODO: Send new OTP email (integrate with email service)

	log.Printf("ResendOTP: new OTP sent for email %s", email)
	return &ResendOTPResponse{
		OTP: newOTP,
	}, nil
}

/*
GetDetailCustomer
mengambil detail customer berdasarkan context
*/
func (s *customerService) GetDetailCustomer(ctx context.Context) (*model.CustomerModel, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	customer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if customer == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "customer"))
	}

	return customer, nil
}

/*
UpdateCustomer
memperbarui data customer dengan validasi field terbatas
*/
func (s *customerService) UpdateCustomer(ctx context.Context, updateData *model.CustomerModel) error {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existingCustomer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existingCustomer == nil {
		return errors.New(fmt.Sprintf(message.NotFound, "customer"))
	}

	existingCustomer.FullName = updateData.FullName
	existingCustomer.PhoneNumber = updateData.PhoneNumber
	existingCustomer.ProfilePhoto = updateData.ProfilePhoto
	existingCustomer.Address = updateData.Address
	existingCustomer.UpdatedAt = time.Now()

	err = s.repo.UpdateCustomer(existingCustomer)
	if err != nil {
		return errors.New(message.InternalError)
	}

	return nil
}

/*
DeleteCustomer
menghapus customer berdasarkan context
*/
func (s *customerService) DeleteCustomer(ctx context.Context) error {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existingCustomer, err := s.repo.GetDetailCustomer(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existingCustomer == nil {
		return errors.New(fmt.Sprintf(message.NotFound, "customer"))
	}

	err = s.repo.DeleteCustomer(id)
	if err != nil {
		return errors.New(message.InternalError)
	}

	return nil
}

/*
UploadIdentity
mengunggah atau memperbarui identitas customer
*/
func (s *customerService) UploadIdentity(ctx context.Context, ktpURL string) error {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existingIdentity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return errors.New(message.InternalError)
	}

	if existingIdentity != nil && existingIdentity.Status == "approved" {
		return errors.New(message.IdentityAlreadyUploaded)
	}

	if existingIdentity != nil {
		existingIdentity.KTPURL = ktpURL
		existingIdentity.Verified = false
		existingIdentity.Status = "pending"
		existingIdentity.Reason = ""
		existingIdentity.VerifiedAt = nil
		existingIdentity.UpdatedAt = time.Now()
		err = s.repo.UpdateIdentity(existingIdentity)
		if err != nil {
			return errors.New(message.InternalError)
		}
		log.Printf("UploadIdentity: updated existing identity for user %s", id)
	} else {
		identityID := uuid.New().String()
		identity := &model.IdentityModel{
			ID:         identityID,
			UserID:     id,
			KTPURL:     ktpURL,
			Verified:   false,
			Status:     "pending",
			Reason:     "",
			VerifiedAt: nil,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = s.repo.CreateIdentity(identity)
		if err != nil {
			if strings.Contains(err.Error(), "identity already approved") {
				return errors.New(message.IdentityAlreadyUploaded)
			}
			return errors.New(message.InternalError)
		}
		log.Printf("UploadIdentity: created new identity for user %s", id)
	}

	return nil
}

/*
UpdateIdentity
memperbarui identitas customer dengan upload KTP baru
*/
func (s *customerService) UpdateIdentity(ctx context.Context, ktpURL string) error {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existingIdentity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existingIdentity == nil {
		return errors.New(fmt.Sprintf(message.NotFound, "identity"))
	}

	existingIdentity.KTPURL = ktpURL
	existingIdentity.Verified = false
	existingIdentity.Status = "pending"
	existingIdentity.Reason = ""
	existingIdentity.VerifiedAt = nil
	existingIdentity.UpdatedAt = time.Now()

	err = s.repo.UpdateIdentity(existingIdentity)
	if err != nil {
		return errors.New(message.InternalError)
	}
	log.Printf("UpdateIdentity: updated identity for user %s", id)
	return nil
}

/*
CheckIdentityExists
memeriksa apakah identitas customer sudah ada
*/
func (s *customerService) CheckIdentityExists(ctx context.Context) error {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	exists, err := s.repo.CheckIdentityExists(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if exists {
		return errors.New(message.IdentityAlreadyUploaded)
	}
	return nil
}

/*
GetIdentityStatus
mengambil status identitas customer
*/
func (s *customerService) GetIdentityStatus(ctx context.Context) (*model.IdentityModel, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	identity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if identity == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "identity"))
	}

	return identity, nil
}

/*
CreateBooking
membuat booking baru dan mengembalikan detail booking
*/
func (s *customerService) CreateBooking(ctx context.Context, req CreateBookingRequest) (*model.BookingDetailDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	// Query identity to set IdentityID
	identity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	totalDays := int(end.Sub(start).Hours() / 24)

	var rental, deposit int
	for _, item := range req.Items {
		rental += item.SubtotalRental
		deposit += item.SubtotalDeposit
	}
	total := rental + deposit - req.Discount
	outstanding := total

	bookingID := uuid.New().String()

	lockedUntil := time.Now().Add(30 * time.Minute)

	booking := &model.BookingModel{
		ID:                   bookingID,
		LockedUntil:          lockedUntil,
		TimeRemainingMinutes: 0,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		TotalDays:            totalDays,
		DeliveryType:         req.DeliveryType,
		Rental:               rental,
		Deposit:              deposit,
		Discount:             req.Discount,
		Total:                total,
		Outstanding:          outstanding,
		UserID:               id,
		IdentityID:           nil, // Will set below
		Status:               "pending",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Set IdentityID if identity exists
	if identity != nil {
		booking.IdentityID = &identity.ID
	}

	// Set hoster_id from the first item if not set
	if booking.HosterID == "" && len(req.Items) > 0 {
		hosterID, err := s.repo.GetHosterIDByItemID(req.Items[0].ItemID)
		if err != nil {
			log.Printf("CreateBooking: error getting hoster_id: %v", err)
			return nil, errors.New(message.InternalError)
		}
		booking.HosterID = hosterID
	}

	items := make([]model.BookingItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.BookingItem{
			ID:              uuid.New().String(),
			BookingID:       bookingID,
			ItemID:          item.ItemID,
			Name:            item.Name,
			Quantity:        item.Quantity,
			PricePerDay:     item.PricePerDay,
			DepositPerUnit:  item.DepositPerUnit,
			SubtotalRental:  item.SubtotalRental,
			SubtotalDeposit: item.SubtotalDeposit,
		}
	}

	customer := model.BookingCustomer{
		ID:        uuid.New().String(),
		BookingID: bookingID,
		Name:      req.Customer.Name,
		Phone:     req.Customer.Phone,
		Email:     req.Customer.Email,
		Address:   req.Customer.Address,
		Notes:     req.Customer.Notes,
	}

	detail, err := s.repo.CreateBooking(booking, items, customer)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

/*
GetBookingsByUserID
mengambil daftar booking berdasarkan user ID
*/
func (s *customerService) GetBookingsByUserID(ctx context.Context) ([]model.BookingListDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	bookings, err := s.repo.GetBookingsByUserID(id)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

/*
GetListBookings
mengambil daftar semua booking berdasarkan user yang sedang login
*/
func (s *customerService) GetListBookings(ctx context.Context) ([]model.BookingListDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	bookings, err := s.repo.GetBookingsByUserID(id)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

/*
GetDetailBooking
mengambil detail booking berdasarkan booking ID dengan validasi kepemilikan
*/
func (s *customerService) GetDetailBooking(ctx context.Context, bookingID string) (*model.BookingDetailDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	bookingDetail, err := s.repo.GetBookingDetail(bookingID)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if bookingDetail == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "booking"))
	}

	// Validasi bahwa booking milik user yang sedang login
	if bookingDetail.Booking.UserID != id {
		return nil, errors.New(message.Unauthorized)
	}

	return bookingDetail, nil
}

/*
CustomerResponse
format response untuk data customer dengan token
*/
type CustomerResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
CreateCustomerResponse
format response setelah pembuatan customer
*/
type CreateCustomerResponse struct {
	Customer *model.CustomerModel `json:"customer"`
	OTP      string               `json:"otp"`
}

/*
CustomerService
interface untuk operasi service customer
*/
type CustomerService interface {
	LoginCustomer(email, password string) (*CustomerResponse, error)
	CreateCustomer(customer *model.CustomerModel) (*CreateCustomerResponse, error)
	GetDetailCustomer(ctx context.Context) (*model.CustomerModel, error)
	UpdateCustomer(ctx context.Context, updateData *model.CustomerModel) error
	DeleteCustomer(ctx context.Context) error
	UploadIdentity(ctx context.Context, ktpURL string) error
	UpdateIdentity(ctx context.Context, ktpURL string) error
	CheckIdentityExists(ctx context.Context) error
	GetIdentityStatus(ctx context.Context) (*model.IdentityModel, error)
	CreateBooking(ctx context.Context, req CreateBookingRequest) (*model.BookingDetailDTO, error)
	GetBookingsByUserID(ctx context.Context) ([]model.BookingListDTO, error)
	GetListBookings(ctx context.Context) ([]model.BookingListDTO, error)
	GetDetailBooking(ctx context.Context, bookingID string) (*model.BookingDetailDTO, error)
	SendOTP(email string, otp string) error
	ResendOTP(email string) (*ResendOTPResponse, error)
}

/*
CreateBookingRequest
berisi data untuk membuat booking baru
*/
type CreateBookingRequest struct {
	StartDate    string                `json:"start_date"`
	EndDate      string                `json:"end_date"`
	DeliveryType string                `json:"delivery_type"`
	Items        []CreateBookingItem   `json:"items"`
	Customer     CreateBookingCustomer `json:"customer"`
	Delivery     int                   `json:"delivery"`
	Discount     int                   `json:"discount"`
}

/*
CreateBookingItem
berisi data item dalam booking
*/
type CreateBookingItem struct {
	ItemID          string `json:"item_id"`
	Name            string `json:"name"`
	Quantity        int    `json:"quantity"`
	PricePerDay     int    `json:"price_per_day"`
	DepositPerUnit  int    `json:"deposit_per_unit"`
	SubtotalRental  int    `json:"subtotal_rental"`
	SubtotalDeposit int    `json:"subtotal_deposit"`
}

/*
CreateBookingCustomer
berisi data customer dalam booking
*/
type CreateBookingCustomer struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"delivery_address"`
	Notes   string `json:"notes"`
}

/*
ResendOTPResponse
format response untuk pengiriman ulang OTP
*/
type ResendOTPResponse struct {
	OTP string `json:"otp"`
}

/*
NewCustomerService
membuat instance CustomerService dengan repository
*/
func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}
