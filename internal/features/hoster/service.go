package hoster

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
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/pkg/message"
)

/*
hosterService menyediakan logika bisnis untuk hoster.
Menggunakan repository untuk akses data.
*/
type hosterService struct {
	repo HosterRepository
}

/*
Methods untuk hosterService menangani operasi bisnis hoster, item, dan terms.
Dipanggil oleh handler untuk validasi dan logika.
*/
func (s *hosterService) generateTokenHoster(userID string) (*HosterResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "hoster",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return &HosterResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *hosterService) LoginHoster(email, password string) (*HosterResponse, error) {
	hoster, err := s.repo.FindByEmailHosterForLogin(email)
	if err != nil || hoster == nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)) != nil {
		return nil, errors.New(message.MsgUnauthorized)
	}

	return s.generateTokenHoster(hoster.ID)
}

func (s *hosterService) CreateHoster(hoster *model.HosterModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(hoster.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	hoster.PasswordHash = string(hash)
	hoster.CreatedAt = time.Now()
	hoster.UpdatedAt = time.Now()

	err = s.repo.CreateHoster(hoster)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgHosterEmailExists)
		}
		return errors.New(message.MsgInternalServerError)
	}

	return nil
}

func (s *hosterService) GetDetailHoster(ctx context.Context) (*model.HosterModel, error) {
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	hoster, err := s.repo.GetDetailHoster(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if hoster == nil {
		return nil, errors.New(message.MsgHosterNotFound)
	}

	return hoster, nil
}

func (s *hosterService) CreateItem(ctx context.Context, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	input.Name = strings.TrimSpace(input.Name)
	// Handle Description sebagai []string: Trim setiap string di slice

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}

	if input.CategoryID == "" {
		return nil, errors.New(message.MsgBadRequest)
	}

	if input.Stock < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	existing, err := s.repo.FindItemNameByUserID(input.Name, userID)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if existing != nil {
		return nil, errors.New(message.MsgBadRequest)
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateItem(input); err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return s.repo.FindItemNameByID(input.ID)
}

func (s *hosterService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}

	item, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if item == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}

	return item, nil
}

func (s *hosterService) GetAllItems() ([]*model.ItemModel, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return items, nil
}

func (s *hosterService) UpdateItem(ctx context.Context, id string, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if existing == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New(message.MsgUnauthorized)
	}

	input.Name = strings.TrimSpace(input.Name)
	// Handle Description sebagai []string: Trim setiap string di slice

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}

	if input.Stock < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.MsgBadRequest)
	}

	input.ID = id
	input.UserID = userID
	input.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(input); err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return s.repo.FindItemNameByID(id)
}

func (s *hosterService) DeleteItem(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if existing == nil {
		return errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return errors.New(message.MsgUnauthorized)
	}

	return s.repo.DeleteItem(id)
}

func (s *hosterService) CreateTermsAndConditions(ctx context.Context, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if input.Description == nil || len(input.Description) == 0 {
		return nil, errors.New(message.MsgTnCDescriptionRequired)
	}

	// Jika ItemID disediakan, cek apakah item milik user (opsional, tambahkan jika diperlukan)
	if input.ItemID != "" {
		item, err := s.repo.FindItemNameByID(input.ItemID)
		if err != nil {
			return nil, errors.New(message.MsgInternalServerError)
		}
		if item == nil || item.UserID != userID {
			return nil, errors.New(message.MsgUnauthorized)
		}
	}

	existing, err := s.repo.FindTermsAndConditionsByUserIDAndDescription(userID, input.Description)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if existing != nil {
		return nil, errors.New(message.MsgBadRequest)
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateTermsAndConditions(input); err != nil {
		// Handle database duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"tnc_user_id_key\"") {
			return nil, errors.New(message.MsgBadRequest)
		}
		return nil, errors.New(message.MsgInternalServerError)
	}

	return s.repo.FindTermsAndConditionsByID(input.ID)
}

func (s *hosterService) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	tac, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if tac == nil {
		return nil, errors.New(message.MsgTnCNotFound)
	}

	return tac, nil
}

func (s *hosterService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	tacs, err := s.repo.GetAllTermsAndConditions()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return tacs, nil
}

func (s *hosterService) UpdateTermsAndConditions(ctx context.Context, id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if existing == nil {
		return nil, errors.New(message.MsgTnCNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New(message.MsgUnauthorized)
	}

	existing.Description = input.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateTermsAndConditions(existing); err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}

	return s.repo.FindTermsAndConditionsByID(id)
}

func (s *hosterService) DeleteTermsAndConditions(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New(message.MsgUnauthorized)
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return errors.New(message.MsgInternalServerError)
	}
	if existing == nil {
		return errors.New(message.MsgTnCNotFound)
	}
	if existing.UserID != userID {
		return errors.New(message.MsgUnauthorized)
	}

	return s.repo.DeleteTermsAndConditions(id)
}

func (s *hosterService) GetIdentityCustomer(ctx context.Context, userID string) (*model.IdentityModel, error) {
	// Ambil ID admin/hoster dari context (untuk otorisasi)
	adminID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New(message.MsgUnauthorized)
	}

	log.Printf("GetIdentityCustomer: adminID from context: %s, userID param: %s", adminID, userID) // Tambahkan log

	// Panggil repository
	identity, err := s.repo.GetIdentityCustomer(userID)
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	if identity == nil {
		return nil, errors.New("Identity not found")
	}

	return identity, nil
}

func (s *hosterService) UpdateIdentityStatus(ctx context.Context, identityID string, status string, rejectedReason string) error {
	// Ambil ID hoster dari context
	hosterID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		log.Printf("UpdateIdentityStatus: unauthorized access")
		return errors.New(message.MsgUnauthorized)
	}

	log.Printf("UpdateIdentityStatus: hoster %s updating identity %s to status %s", hosterID, identityID, status)

	// Validasi status
	if status != "approved" && status != "rejected" {
		log.Printf("UpdateIdentityStatus: invalid status %s", status)
		return errors.New("Invalid status")
	}

	// Cek apakah identity ada
	identity, err := s.repo.GetIdentityCustomerByID(identityID)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error getting identity %s: %v", identityID, err)
		return errors.New(message.MsgInternalServerError)
	}
	if identity == nil {
		log.Printf("UpdateIdentityStatus: identity %s not found", identityID)
		return errors.New("Identity not found")
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

	// Update status
	err = s.repo.UpdateIdentityStatus(identityID, status, rejectedReason, verified, verifiedAt)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error updating identity %s: %v", identityID, err)
		return errors.New(message.MsgInternalServerError)
	}

	// Jika approved atau rejected, update booking identity status
	if status == "approved" || status == "rejected" {
		log.Printf("UpdateIdentityStatus: updating booking identity for user %s to %s", identity.UserID, status)
		err = s.repo.UpdateBookingIdentityStatusByUserID(identity.UserID, status)
		if err != nil {
			log.Printf("UpdateIdentityStatus: error updating booking identity status: %v", err)
			return errors.New(message.MsgInternalServerError)
		}
	}

	log.Printf("UpdateIdentityStatus: successfully updated identity %s to status %s", identityID, status)
	return nil
}

func (s *hosterService) VerifyIdentity(ctx context.Context, identityID string, status string, rejectionReason string) error {
	// Ambil ID hoster dari context
	hosterID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		log.Printf("VerifyIdentity: unauthorized access")
		return errors.New(message.MsgUnauthorized)
	}

	log.Printf("VerifyIdentity: hoster %s verifying identity %s with status %s", hosterID, identityID, status)

	// Validasi status
	if status != "approved" && status != "rejected" {
		log.Printf("VerifyIdentity: invalid status %s", status)
		return errors.New("Invalid status")
	}

	// Cek apakah identity ada
	identity, err := s.repo.GetIdentityCustomerByID(identityID)
	if err != nil {
		log.Printf("VerifyIdentity: error getting identity %s: %v", identityID, err)
		return errors.New(message.MsgInternalServerError)
	}
	if identity == nil {
		log.Printf("VerifyIdentity: identity %s not found", identityID)
		return errors.New("Identity not found")
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

	// Update status
	err = s.repo.UpdateIdentityStatus(identityID, status, rejectionReason, verified, verifiedAt)
	if err != nil {
		log.Printf("VerifyIdentity: error updating identity %s: %v", identityID, err)
		return errors.New(message.MsgInternalServerError)
	}

	// Jika approved, update booking identity status
	if status == "approved" {
		log.Printf("VerifyIdentity: updating booking identity for user %s", identity.UserID)
		err = s.repo.UpdateBookingIdentityStatusByUserID(identity.UserID, "approved")
		if err != nil {
			log.Printf("VerifyIdentity: error updating booking identity status: %v", err)
			return errors.New(message.MsgInternalServerError)
		}
	}

	log.Printf("VerifyIdentity: successfully updated identity %s to status %s", identityID, status)
	return nil
}

/*
HosterResponse berisi data respons autentikasi hoster.
Digunakan untuk mengembalikan token dan info user.
*/
type HosterResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
HosterService mendefinisikan kontrak operasi bisnis hoster.
Diimplementasikan oleh hosterService.
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
	UpdateIdentityStatus(ctx context.Context, identityID string, status string, rejectedReason string) error // Tambahkan ini
}

/*
NewHosterService membuat instance HosterService.
Menginisialisasi service dengan repository.
*/
func NewHosterService(repo HosterRepository) HosterService {
	return &hosterService{repo: repo}
}
