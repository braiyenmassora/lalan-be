package tnc

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
TnCService adalah kontrak untuk logika bisnis T&C.
*/
type TnCService interface {
	CreateTnC(hosterID string, req *dto.CreateTnCRequest) (*dto.TnCResponse, error)
	UpdateTnC(hosterID, tncID string, req *dto.UpdateTnCRequest) (*dto.TnCResponse, error)
	GetTnC(hosterID string) (*dto.TnCResponse, error)
}

/*
tncService adalah implementasi service untuk T&C.
*/
type tncService struct {
	repo TnCRepository
}

/*
NewTnCService membuat instance service dengan dependency injection.

Output:
- TnCService siap digunakan
*/
func NewTnCService(repo TnCRepository) TnCService {
	return &tncService{repo: repo}
}

/*
CreateTnC membuat T&C baru untuk hoster.

Alur kerja:
1. Validasi userID dan description tidak kosong
2. Build domain entity
3. Panggil repository
4. Return response

Output sukses:
- (*dto.TnCResponse, nil)
Output error:
- (nil, error) → unauthorized / bad request / internal error
*/
func (s *tncService) CreateTnC(hosterID string, req *dto.CreateTnCRequest) (*dto.TnCResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	if req == nil || len(req.Description) == 0 {
		return nil, errors.New(message.BadRequest)
	}

	// Build entity
	tnc := &domain.TermsAndConditions{
		ID:          uuid.New().String(),
		UserID:      hosterID,
		Description: req.Description,
	}

	// Save to DB
	if err := s.repo.CreateTnC(tnc); err != nil {
		log.Printf("CreateTnC service: repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	// Get created TnC
	created, err := s.repo.GetTnCByHosterID(hosterID)
	if err != nil {
		log.Printf("CreateTnC service: failed to get created tnc for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return created, nil
}

/*
UpdateTnC memperbarui T&C hoster.

Alur kerja:
1. Validasi userID, tncID, dan description tidak kosong
2. Panggil repository update
3. Get updated data
4. Return response

Output sukses:
- (*dto.TnCResponse, nil)
Output error:
- (nil, error) → unauthorized / bad request / not found / internal error
*/
func (s *tncService) UpdateTnC(hosterID, tncID string, req *dto.UpdateTnCRequest) (*dto.TnCResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	if tncID == "" || req == nil || len(req.Description) == 0 {
		return nil, errors.New(message.BadRequest)
	}

	// Update in DB
	if err := s.repo.UpdateTnC(tncID, hosterID, req.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(message.TnCNotFound)
		}
		log.Printf("UpdateTnC service: repo error for hoster %s tnc %s: %v", hosterID, tncID, err)
		return nil, errors.New(message.InternalError)
	}

	// Get updated TnC
	updated, err := s.repo.GetTnCByHosterID(hosterID)
	if err != nil {
		log.Printf("UpdateTnC service: failed to get updated tnc for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return updated, nil
}

/*
GetTnC mengambil T&C berdasarkan hosterID.

Alur kerja:
1. Validasi hosterID tidak kosong
2. Panggil repository
3. Return response

Output sukses:
- (*dto.TnCResponse, nil)
Output error:
- (nil, error) → unauthorized / not found / internal error
*/
func (s *tncService) GetTnC(hosterID string) (*dto.TnCResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	tnc, err := s.repo.GetTnCByHosterID(hosterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(message.TnCNotFound)
		}
		log.Printf("GetTnC service: repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return tnc, nil
}
