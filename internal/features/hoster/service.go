package hoster

import (
	"context"
	"errors"
	"fmt"
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
type hosterService struct
menyediakan logika bisnis untuk operasi hoster
*/
type hosterService struct {
	repo HosterRepository
}

/*
HosterService
mendefinisikan kontrak operasi bisnis hoster
*/
type HosterService interface {
	CreateHoster(*model.HosterModel) error
	LoginHoster(email, password string) (*HosterResponse, error)
	GetDetailHoster(ctx context.Context) (*model.HosterModel, error)
	CreateItem(ctx context.Context, input *model.ItemModel) (*model.ItemModel, error)
	GetItemByID(id string) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	UpdateItem(ctx context.Context, id string, input *model.ItemModel) (*model.ItemModel, error)
	DeleteItem(ctx context.Context, id string) error
	CreateTermsAndConditions(ctx context.Context, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
	UpdateTermsAndConditions(ctx context.Context, id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	DeleteTermsAndConditions(ctx context.Context, id string) error
	GetIdentityCustomer(ctx context.Context, userID string) (*model.IdentityModel, error)
	GetListBookingsCustomer(ctx context.Context, limit int, offset int) ([]model.BookingListDTOHoster, error)
	GetListBookingsCustomerByBookingID(ctx context.Context, bookingID string, limit int, offset int) ([]model.BookingDetailDTOHoster, error)
	GetListCustomer(ctx context.Context) ([]model.CustomerIdentityDTO, error)
	GetDetailCustomer(ctx context.Context, customerID string) (*model.CustomerIdentityDetailDTO, error)
}

/*
NewHosterService
membuat instance baru HosterService dengan repository yang diberikan
*/
func NewHosterService(repo HosterRepository) HosterService {
	return &hosterService{repo: repo}
}

/*
generateTokenHoster
menghasilkan JWT token untuk hoster
*/
// di hoster/service.go — ganti fungsi generateTokenHoster kamu
func (s *hosterService) generateTokenHoster(userID string) (*HosterResponse, error) {
	// Access Token — 15 menit
	accessToken, err := s.generateJWT(userID, "hoster", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	// Refresh Token — 30 hari (JWT juga, biar bisa verify tanpa DB)
	refreshToken, err := s.generateJWT(userID, "hoster", 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &HosterResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 menit
		TokenType:    "Bearer",
		Role:         "hoster",
	}, nil
}

// Helper generate JWT (pakai secret dari env)
func (s *hosterService) generateJWT(userID, role string, expires time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(expires).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTSecret())
}

/*
LoginHoster
memvalidasi login hoster dan menghasilkan token
*/
func (s *hosterService) LoginHoster(email, password string) (*HosterResponse, error) {
	hoster, err := s.repo.FindByEmailHosterForLogin(email)
	if err != nil || hoster == nil {
		return nil, errors.New(message.LoginFailed)
	}

	if bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.LoginFailed)
	}

	return s.generateTokenHoster(hoster.ID)
}

/*
CreateHoster
membuat hoster baru dengan hash password
*/
func (s *hosterService) CreateHoster(hoster *model.HosterModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(hoster.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.InternalError)
	}
	hoster.PasswordHash = string(hash)
	hoster.CreatedAt = time.Now()
	hoster.UpdatedAt = time.Now()

	err = s.repo.CreateHoster(hoster)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(fmt.Sprintf(message.AlreadyExists, "hoster email"))
		}
		return errors.New(message.InternalError)
	}

	return nil
}

/*
GetDetailHoster
mengambil detail hoster berdasarkan context
*/
func (s *hosterService) GetDetailHoster(ctx context.Context) (*model.HosterModel, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	hoster, err := s.repo.GetDetailHoster(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if hoster == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "hoster"))
	}

	return hoster, nil
}

/*
CreateItem
membuat item baru dengan validasi input
*/
func (s *hosterService) CreateItem(ctx context.Context, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	input.Name = strings.TrimSpace(input.Name)
	// Handle Description sebagai []string: Trim setiap string di slice

	if input.Name == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "item name"))
	}

	if input.CategoryID == "" {
		return nil, errors.New(message.BadRequest)
	}

	if input.Stock < 0 {
		return nil, errors.New(message.BadRequest)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.BadRequest)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.BadRequest)
	}

	existing, err := s.repo.FindItemNameByUserID(input.Name, userID)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if existing != nil {
		return nil, errors.New(fmt.Sprintf(message.AlreadyExists, "item"))
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateItem(input); err != nil {
		return nil, errors.New(message.InternalError)
	}

	return s.repo.FindItemNameByID(input.ID)
}

/*
GetItemByID
mengambil item berdasarkan ID
*/
func (s *hosterService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "item ID"))
	}

	item, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if item == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "item"))
	}

	return item, nil
}

/*
GetAllItems
mengambil semua item
*/
func (s *hosterService) GetAllItems() ([]*model.ItemModel, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return items, nil
}

/*
UpdateItem
memperbarui item berdasarkan ID dengan validasi kepemilikan
*/
func (s *hosterService) UpdateItem(ctx context.Context, id string, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if existing == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "item"))
	}
	if existing.UserID != userID {
		return nil, errors.New(message.Unauthorized)
	}

	input.Name = strings.TrimSpace(input.Name)
	// Handle Description sebagai []string: Trim setiap string di slice

	if input.Name == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "item name"))
	}

	if input.Stock < 0 {
		return nil, errors.New(message.BadRequest)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.BadRequest)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.BadRequest)
	}

	input.ID = id
	input.UserID = userID
	input.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(input); err != nil {
		return nil, errors.New(message.InternalError)
	}

	return s.repo.FindItemNameByID(id)
}

/*
DeleteItem
menghapus item berdasarkan ID dengan validasi kepemilikan
*/
func (s *hosterService) DeleteItem(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existing == nil {
		return errors.New(fmt.Sprintf(message.NotFound, "item"))
	}
	if existing.UserID != userID {
		return errors.New(message.Unauthorized) // Diperbaiki: Hapus 'nil,' agar hanya return error
	}

	return s.repo.DeleteItem(id)
}

/*
CreateTermsAndConditions
membuat syarat dan ketentuan baru dengan validasi input
*/
func (s *hosterService) CreateTermsAndConditions(ctx context.Context, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	if input.Description == nil || len(input.Description) == 0 {
		return nil, errors.New(fmt.Sprintf(message.Required, "terms and conditions description"))
	}

	// Jika ItemID disediakan, cek apakah item milik user (opsional, tambahkan jika diperlukan)
	if input.ItemID != "" {
		item, err := s.repo.FindItemNameByID(input.ItemID)
		if err != nil {
			return nil, errors.New(message.InternalError)
		}
		if item == nil || item.UserID != userID {
			return nil, errors.New(message.Unauthorized)
		}
	}

	existing, err := s.repo.FindTermsAndConditionsByUserIDAndDescription(userID, input.Description)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if existing != nil {
		return nil, errors.New(fmt.Sprintf(message.AlreadyExists, "terms and conditions"))
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateTermsAndConditions(input); err != nil {
		// Handle database duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"tnc_user_id_key\"") {
			return nil, errors.New(fmt.Sprintf(message.AlreadyExists, "terms and conditions"))
		}
		return nil, errors.New(message.InternalError)
	}

	return s.repo.FindTermsAndConditionsByID(input.ID)
}

/*
FindTermsAndConditionsByID
mencari syarat dan ketentuan berdasarkan ID
*/
func (s *hosterService) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	tac, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if tac == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "terms and conditions"))
	}

	return tac, nil
}

/*
GetAllTermsAndConditions
mengambil semua syarat dan ketentuan
*/
func (s *hosterService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	tacs, err := s.repo.GetAllTermsAndConditions()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return tacs, nil
}

/*
UpdateTermsAndConditions
memperbarui syarat dan ketentuan berdasarkan ID dengan validasi kepemilikan
*/
func (s *hosterService) UpdateTermsAndConditions(ctx context.Context, id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if existing == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "terms and conditions"))
	}
	if existing.UserID != userID {
		return nil, errors.New(message.Unauthorized)
	}

	existing.Description = input.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateTermsAndConditions(existing); err != nil {
		return nil, errors.New(message.InternalError)
	}

	return s.repo.FindTermsAndConditionsByID(id)
}

/*
DeleteTermsAndConditions
menghapus syarat dan ketentuan berdasarkan ID dengan validasi kepemilikan
*/
func (s *hosterService) DeleteTermsAndConditions(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.Unauthorized)
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return errors.New(message.InternalError)
	}
	if existing == nil {
		return errors.New(fmt.Sprintf(message.NotFound, "terms and conditions"))
	}
	if existing.UserID != userID {
		return errors.New(message.Unauthorized) // Diperbaiki: Hapus 'nil,' agar hanya return error
	}

	return s.repo.DeleteTermsAndConditions(id)
}

/*
GetIdentityCustomer
mengambil identitas customer berdasarkan userID
*/
func (s *hosterService) GetIdentityCustomer(ctx context.Context, userID string) (*model.IdentityModel, error) {
	if userID == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "user ID"))
	}

	identity, err := s.repo.GetIdentityCustomer(userID)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if identity == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "identity"))
	}

	return identity, nil
}

/*
GetListBookingsForHoster
mengambil daftar booking yang dimiliki oleh hoster berdasarkan context dengan limit dan offset
*/
func (s *hosterService) GetListBookingsForHoster(ctx context.Context, limit int, offset int) ([]model.BookingListDTOHoster, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	bookings, err := s.repo.GetListBookingsCustomer(userID, limit, offset)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return bookings, nil
}

/*
GetListBookingsCustomer
mengambil daftar booking yang dimiliki oleh hoster berdasarkan context dengan limit dan offset
*/
func (s *hosterService) GetListBookingsCustomer(ctx context.Context, limit int, offset int) ([]model.BookingListDTOHoster, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	bookings, err := s.repo.GetListBookingsCustomer(userID, limit, offset)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return bookings, nil
}

/*
GetListBookingsCustomerByBookingID
mengambil daftar booking yang dimiliki oleh hoster berdasarkan booking ID dengan limit dan offset
*/
func (s *hosterService) GetListBookingsCustomerByBookingID(ctx context.Context, bookingID string, limit int, offset int) ([]model.BookingDetailDTOHoster, error) {
	hosterID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	if bookingID == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "booking ID"))
	}

	bookings, err := s.repo.GetListBookingsCustomerByBookingID(hosterID, bookingID, limit, offset)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return bookings, nil
}

/*
GetListCustomer
mengambil daftar customer unik yang telah booking dengan hoster tertentu berdasarkan context
*/
func (s *hosterService) GetListCustomer(ctx context.Context) ([]model.CustomerIdentityDTO, error) {
	hosterID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	customers, err := s.repo.GetListCustomer(hosterID)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	return customers, nil
}

/*
GetDetailCustomer
mengambil detail customer berdasarkan customer ID dengan validasi hoster
*/
func (s *hosterService) GetDetailCustomer(ctx context.Context, customerID string) (*model.CustomerIdentityDetailDTO, error) {
	hosterID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.Unauthorized)
	}

	if customerID == "" {
		return nil, errors.New(fmt.Sprintf(message.Required, "customer ID"))
	}

	customer, err := s.repo.GetDetailCustomer(customerID, hosterID)
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	if customer == nil {
		return nil, errors.New(fmt.Sprintf(message.NotFound, "customer"))
	}

	return customer, nil
}

/*
HosterResponse
struct untuk response login hoster
*/
type HosterResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Role         string `json:"role"`
}
