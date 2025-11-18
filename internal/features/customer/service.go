package customer

import (
	"context"
	"errors"
	"lalan-be/internal/config"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/pkg/message"
	"log"
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
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	// Cek jika identity sudah ada
	existingIdentity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}

	if existingIdentity != nil {
		// Update existing: replace KTP, reset status
		existingIdentity.KTPURL = ktpURL
		existingIdentity.Verified = false
		existingIdentity.Status = "pending"
		existingIdentity.RejectedReason = ""
		existingIdentity.VerifiedAt = nil
		existingIdentity.UpdatedAt = time.Now()
		err = s.repo.UpdateIdentity(existingIdentity)
		if err != nil {
			return errors.New(message.MsgInternalServerError)
		}
		log.Printf("UploadIdentity: updated existing identity for user %s", id)
	} else {
		// Create new identity
		identityID := uuid.New().String()
		identity := &model.IdentityModel{
			ID:             identityID,
			UserID:         id,
			KTPURL:         ktpURL,
			Verified:       false,
			Status:         "pending",
			RejectedReason: "",
			VerifiedAt:     nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err = s.repo.CreateIdentity(identity)
		if err != nil {
			return errors.New(message.MsgInternalServerError)
		}
		log.Printf("UploadIdentity: created new identity for user %s", id)
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

func (s *customerService) CreateBooking(ctx context.Context, req CreateBookingRequest) (*model.BookingDetailDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	// Cek identity customer
	identity, err := s.repo.GetIdentityByUserID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if identity == nil {
		return nil, errors.New("Identity not found")
	}

	// Hitung total_days
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	totalDays := int(end.Sub(start).Hours() / 24)

	// Hitung price
	var rental, deposit int
	for _, item := range req.Items {
		rental += item.SubtotalRental
		deposit += item.SubtotalDeposit
	}
	total := rental + deposit + req.Delivery - req.Discount
	outstanding := total

	// Generate ID dan code
	bookingID := uuid.New().String()
	code := "BK" + time.Now().Format("060102") + uuid.New().String()[:4]

	// Locked until (30 menit)
	lockedUntil := time.Now().Add(30 * time.Minute)

	// Build models dengan field flat
	booking := &model.BookingModel{
		ID:                   bookingID,
		Code:                 code,
		LockedUntil:          lockedUntil,
		TimeRemainingMinutes: 0, // Akan dihitung nanti
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		TotalDays:            totalDays,
		DeliveryType:         req.DeliveryType,
		Rental:               rental,
		Deposit:              deposit,
		Delivery:             req.Delivery,
		Discount:             req.Discount,
		Total:                total,
		Outstanding:          outstanding,
		UserID:               id,
		IdentityID:           &identity.ID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
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
		ID:              uuid.New().String(),
		BookingID:       bookingID,
		Name:            req.Customer.Name,
		Phone:           req.Customer.Phone,
		Email:           req.Customer.Email,
		DeliveryAddress: req.Customer.DeliveryAddress,
		Notes:           req.Customer.Notes,
	}

	bookingIdentity := model.BookingIdentity{
		ID:              uuid.New().String(),
		BookingID:       bookingID,
		Uploaded:        true,
		Status:          identity.Status,
		RejectionReason: &identity.RejectedReason,
		ReuploadAllowed: identity.Status == "rejected",
		EstimatedTime:   "Maksimal 30 menit",
		StatusCheckURL:  "/api/v1/customer/identity-status",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Insert
	err = s.repo.CreateBooking(booking, items, customer, bookingIdentity)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	// Hitung time_remaining_minutes
	timeRemaining := int(lockedUntil.Sub(time.Now()).Minutes())

	// Build DTO
	dto := &model.BookingDetailDTO{
		Booking:  *booking,
		Items:    items,
		Customer: customer,
		Identity: bookingIdentity,
	}
	dto.Booking.TimeRemainingMinutes = timeRemaining

	return dto, nil
}

func (s *customerService) GetBookingsByUserID(ctx context.Context) ([]model.BookingListDTO, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	bookings, err := s.repo.GetBookingsByUserID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return bookings, nil
}

func (s *customerService) GetListBookings(ctx context.Context) ([]model.BookingListDTO, error) {
	// Opsional: Cek role admin/hoster
	bookings, err := s.repo.GetListBookings()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return bookings, nil
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
	GetIdentityStatus(ctx context.Context) (*model.IdentityModel, error)
	CreateBooking(ctx context.Context, req CreateBookingRequest) (*model.BookingDetailDTO, error) // Tambahkan ini
	GetBookingsByUserID(ctx context.Context) ([]model.BookingListDTO, error)
	GetListBookings(ctx context.Context) ([]model.BookingListDTO, error) // Rename
}

// Tambahkan struct request
type CreateBookingRequest struct {
	StartDate    string                `json:"start_date"`
	EndDate      string                `json:"end_date"`
	DeliveryType string                `json:"delivery_type"`
	Items        []CreateBookingItem   `json:"items"`
	Customer     CreateBookingCustomer `json:"customer"`
	Delivery     int                   `json:"delivery"`
	Discount     int                   `json:"discount"`
}

type CreateBookingItem struct {
	ItemID          string `json:"item_id"`
	Name            string `json:"name"`
	Quantity        int    `json:"quantity"`
	PricePerDay     int    `json:"price_per_day"`
	DepositPerUnit  int    `json:"deposit_per_unit"`
	SubtotalRental  int    `json:"subtotal_rental"`
	SubtotalDeposit int    `json:"subtotal_deposit"`
}

type CreateBookingCustomer struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	DeliveryAddress string `json:"delivery_address"`
	Notes           string `json:"notes"`
}

/*
NewCustomerService membuat instance CustomerService.
Menginisialisasi service dengan repository.
*/
func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}
