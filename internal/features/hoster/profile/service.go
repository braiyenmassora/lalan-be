package profile

import (
	"database/sql"
	"errors"
	"log"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
HosterProfileService adalah kontrak untuk logika bisnis profile hoster.
*/
type HosterProfileService interface {
	GetProfile(hosterID string) (*dto.HosterProfileResponse, error)
	UpdateProfile(hosterID string, req *dto.UpdateHosterProfileRequest) (*dto.HosterProfileResponse, error)
}

/*
hosterProfileService adalah implementasi service untuk profile hoster.
*/
type hosterProfileService struct {
	repo HosterProfileRepository
}

/*
NewHosterProfileService membuat instance service dengan dependency injection.

Output:
- HosterProfileService siap digunakan
*/
func NewHosterProfileService(repo HosterProfileRepository) HosterProfileService {
	return &hosterProfileService{repo: repo}
}

/*
GetProfile mengambil profil hoster.

Alur kerja:
1. Validasi hosterID tidak kosong
2. Panggil repository
3. Return response

Output sukses:
- (*dto.HosterProfileResponse, nil)
Output error:
- (nil, error) → unauthorized / not found / internal error
*/
func (s *hosterProfileService) GetProfile(hosterID string) (*dto.HosterProfileResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	profile, err := s.repo.GetProfile(hosterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(message.ProfileNotFound)
		}
		log.Printf("GetProfile service: repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return profile, nil
}

/*
UpdateProfile memperbarui profil hoster (address dan phone_number).

Alur kerja:
1. Validasi hosterID dan request tidak kosong
2. Panggil repository update
3. Get updated profile
4. Return response

Output sukses:
- (*dto.HosterProfileResponse, nil)
Output error:
- (nil, error) → unauthorized / bad request / not found / internal error
*/
func (s *hosterProfileService) UpdateProfile(hosterID string, req *dto.UpdateHosterProfileRequest) (*dto.HosterProfileResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	if req == nil {
		return nil, errors.New(message.BadRequest)
	}

	// Update in DB
	if err := s.repo.UpdateProfile(hosterID, req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(message.ProfileNotFound)
		}
		log.Printf("UpdateProfile service: repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	// Get updated profile
	updated, err := s.repo.GetProfile(hosterID)
	if err != nil {
		log.Printf("UpdateProfile service: failed to get updated profile for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return updated, nil
}
