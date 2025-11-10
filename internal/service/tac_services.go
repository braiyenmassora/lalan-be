package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"

	"github.com/google/uuid"
)

/*
Mendefinisikan operasi service TAC.
Menyediakan method untuk menambah TAC dengan hasil sukses atau error.
*/
type TermsAndConditionsService interface {
	AddTermsAndConditions(input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
}

/*
Implementasi service TAC dengan repository.
*/
type termsAndConditionsService struct {
	repo repository.TermsAndConditionsRepository
}

/*
Membuat service TAC.
Mengembalikan instance TermsAndConditionsService yang siap digunakan.
*/
func NewTermsAndConditionsService(repo repository.TermsAndConditionsRepository) TermsAndConditionsService {
	return &termsAndConditionsService{repo: repo}
}

/*
Menambahkan TAC baru.
Mengembalikan data TAC yang dibuat atau error jika validasi/gagal.
*/
func (s *termsAndConditionsService) AddTermsAndConditions(input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	// Validasi input
	if len(input.Description) == 0 {
		return nil, errors.New(message.MsgTermAndConditionsDescriptionRequired)
	}

	if input.UserID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}

	// Cek apakah user sudah punya TAC
	existing, err := s.repo.FindByUserID(input.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgTermAndConditionsAlreadyExists)
	}

	// Menghasilkan ID unik
	input.ID = uuid.New().String()

	// Menyimpan TAC ke database
	if err := s.repo.CreateTermAndConditions(input); err != nil {
		return nil, err
	}

	// Mendapatkan TAC yang baru dibuat
	return s.repo.FindByUserID(input.UserID)
}
