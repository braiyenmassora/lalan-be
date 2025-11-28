// This file contains the IdentityService implementation for handling KTP verification.
package identity

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/utils"
)

/*
IdentityRepo adalah kontrak (interface) untuk operasi database pada tabel identity.
Implementasi nyata berada di IdentityRepository (struct dengan sqlx).
*/
type IdentityRepo interface {
	UploadKTP(req *dto.UploadIdentityByCustomerRequest) error
	GetStatusKTP(userID string) (*domain.Identity, error)
}

/*
IdentityService adalah layer business logic untuk fitur verifikasi KTP.
Bertanggung jawab atas:
• Orkestrasi antara storage (upload/delete file) dan repository (DB)
• Menjaga konsistensi data (file terhapus saat record dihapus)
• Validasi dasar userID dan file
Tidak boleh ada detail HTTP atau query SQL di sini.
*/
type IdentityService struct {
	repo    IdentityRepo
	storage utils.Storage
}

/*
NewIdentityService membuat instance service yang siap digunakan.
Dependency injection untuk repo dan storage memudahkan unit testing dan pergantian implementasi.

Output:
- *IdentityService yang sudah terkoneksi ke repository dan storage.
*/
func NewIdentityService(repo IdentityRepo, storage utils.Storage) *IdentityService {
	return &IdentityService{
		repo:    repo,
		storage: storage,
	}
}

/*
UploadKTP menangani upload KTP pertama kali oleh customer.

Alur kerja:
1. Validasi userID tidak kosong
2. Hapus semua record lama + file di storage (cleanup)
3. Upload file baru ke storage dengan path deterministik: ktp/{userID}.jpg
4. Simpan record baru ke DB dengan status "pending"

Output sukses:
- error = nil → upload berhasil, record tersimpan
Output error:
- error → userID kosong / gagal hapus lama / gagal upload file / gagal insert DB
*/
func (s *IdentityService) UploadKTP(ctx context.Context, userID string, file io.Reader) error {
	if userID == "" {
		return fmt.Errorf("userID required")
	}

	// Kami tidak menghapus record lama pada upload/re-upload.
	// Setiap upload harus dianggap sebagai entri baru (insert) supaya
	// referensi lama masih aman (misalnya booking yang mengacu ke identity lama).

	// Upload file baru — gunakan nama file unik agar tidak menimpa file lama
	// di object storage.
	path := fmt.Sprintf("ktp/%s_%d.jpg", userID, time.Now().UnixNano())
	newURL, err := s.storage.Upload(ctx, file, path, "image/jpeg")
	if err != nil {
		return fmt.Errorf("failed to upload new ktp: %w", err)
	}

	// 3. Simpan record baru (status pending)
	req := &dto.UploadIdentityByCustomerRequest{
		UserID: userID,
		KTPURL: newURL,
	}
	if err := s.repo.UploadKTP(req); err != nil {
		// Jika DB gagal, coba hapus file yang baru di-upload
		_ = s.storage.Delete(ctx, newURL)
		return fmt.Errorf("failed to save identity record: %w", err)
	}

	return nil
}

/*
UpdateKTP menangani re-upload KTP oleh customer yang sudah memiliki record.

Alur kerja:
1. Validasi userID dan file tidak kosong
2. Upload file baru ke storage (overwrite path yang sama)
3. Update record di DB → reset status ke "pending", verified = false

Output sukses:
- error = nil → file terganti, status di-reset ke pending
Output error:
- error → userID/file kosong / gagal upload / gagal update DB (file tetap dihapus jika DB gagal)
*/
func (s *IdentityService) UpdateKTP(ctx context.Context, userID string, file io.Reader) error {
	if userID == "" {
		return fmt.Errorf("userID required")
	}
	if file == nil {
		return fmt.Errorf("file is required")
	}

	// Treat re-upload as new entry → insert new record and keep old ones.
	// Create a unique path (timestamp based) so file objects are not overwritten.
	path := fmt.Sprintf("ktp/%s_%d.jpg", userID, time.Now().UnixNano())
	newURL, err := s.storage.Upload(ctx, file, path, "image/jpeg")
	if err != nil {
		return fmt.Errorf("failed to upload ktp: %w", err)
	}

	// Save as a NEW identity record (insert) — do NOT update or delete old records.
	req := &dto.UploadIdentityByCustomerRequest{
		UserID: userID,
		KTPURL: newURL,
	}
	if err := s.repo.UploadKTP(req); err != nil {
		// Jika DB gagal, coba hapus file yang baru di-upload
		_ = s.storage.Delete(ctx, newURL)
		return fmt.Errorf("failed to update identity record: %w", err)
	}

	return nil
}

/*
GetStatusKTP mengembalikan status verifikasi KTP milik user yang login.

Alur kerja:
1. Validasi userID tidak kosong
2. Ambil record dari repository
3. Mapping ke DTO response

Output sukses:
- (*dto.IdentityStatusDTO, nil) → record ditemukan
- (nil, nil)                → user belum pernah upload KTP
Output error:
- (nil, error)             → userID kosong / error repository
*/
func (s *IdentityService) GetStatusKTP(ctx context.Context, userID string) (*dto.IdentityStatusByCustomerResponse, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	model, err := s.repo.GetStatusKTP(userID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if model == nil {
		return nil, nil
	}

	return &dto.IdentityStatusByCustomerResponse{
		KTPID:      model.ID,
		UserID:     model.UserID,
		KTPURL:     model.KTPURL,
		CreatedAt:  model.CreatedAt,
		Status:     model.Status,
		Verified:   model.Verified,
		Reason:     model.Reason,
		VerifiedAt: model.VerifiedAt,
	}, nil
}
