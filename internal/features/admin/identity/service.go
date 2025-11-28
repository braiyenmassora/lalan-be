package identity

import (
	"fmt"

	"lalan-be/internal/domain"
)

/*
AdminIdentityService mengatur logika bisnis untuk fitur verifikasi identitas oleh admin.
Berperan sebagai lapisan validasi dan koordinasi antara handler dan repository.
*/
type AdminIdentityService struct {
	repo *AdminIdentityRepository
}

/*
NewAdminIdentityService membuat instance service dengan dependency injection.

Output:
- *AdminIdentityService siap digunakan
*/
func NewAdminIdentityService(repo *AdminIdentityRepository) *AdminIdentityService {
	return &AdminIdentityService{repo: repo}
}

/*
GetPendingIdentities mengambil semua identitas yang berstatus 'pending' untuk ditinjau admin.

Alur kerja:
1. Delegasikan ke repository

Output sukses:
- ([]*model.IdentityModel, nil)
Output error:
- (nil, error) → query gagal / DB error
*/
func (s *AdminIdentityService) GetPendingIdentities() ([]*domain.Identity, error) {
	return s.repo.GetPendingIdentities()
}

/*
ValidateIdentity memproses persetujuan atau penolakan KTP oleh admin.

Alur kerja:
1. Validasi status hanya boleh "approved" atau "rejected"
2. Jika rejected → reason wajib diisi
3. Panggil repository untuk update status

Output sukses:
- nil → status berhasil diperbarui
Output error:
- error → status tidak valid / reason kosong saat rejected
*/
func (s *AdminIdentityService) ValidateIdentity(identityID, status, reason string) error {
	if status != "approved" && status != "rejected" {
		return fmt.Errorf("invalid status: must be 'approved' or 'rejected'")
	}

	if status == "rejected" && reason == "" {
		return fmt.Errorf("reason is required when rejecting identity")
	}

	// Use the repository ValidateIdentity which targets a specific identity record id
	// to avoid ambiguity and allow validating historical uploads separately.
	return s.repo.ValidateIdentity(identityID, status, reason)
}

/*
GetIdentity mengambil detail identitas satu user berdasarkan userID.

Alur kerja:
1. Delegasikan ke repository

Output sukses:
- (*model.IdentityModel, nil)
Output error:
- (nil, error) → identitas tidak ditemukan / DB error
*/
func (s *AdminIdentityService) GetIdentity(userID string) (*domain.Identity, error) {
	return s.repo.GetIdentityByUserID(userID)
}
