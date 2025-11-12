package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"

	"github.com/google/uuid"
)

/*
Implementasi service TAC dengan repository.
*/
type termsAndConditionsService struct {
	repo repository.TermsAndConditionsRepository
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

/*
Update TAC berdasarkan ID.
Mengembalikan data TAC yang diupdate atau error jika validasi/gagal.
*/
func (s *termsAndConditionsService) UpdateTermAndConditions(id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	// Validasi ID
	if id == "" {
		return nil, errors.New(message.MsgTermAndConditionsIDRequired)
	}

	// Validasi input
	if len(input.Description) == 0 {
		return nil, errors.New(message.MsgTermAndConditionsDescriptionRequired)
	}

	// Update TAC
	input.ID = id
	if err := s.repo.UpdateTermAndConditions(input); err != nil {
		return nil, err
	}

	// Mendapatkan TAC yang sudah diupdate
	return s.repo.FindTermAndConditionsByID(id)
}

/*
Update TAC berdasarkan userID (hanya description).
Mengembalikan data TAC yang diupdate atau error jika validasi/gagal.
*/
func (s *termsAndConditionsService) UpdateTermsAndConditionsByUserID(userID string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	// Validasi input
	if len(input.Description) == 0 {
		return nil, errors.New(message.MsgTermAndConditionsDescriptionRequired)
	}

	// Cari TAC existing berdasarkan userID
	existing, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgTermAndConditionsNotFound)
	}

	// Update description
	existing.Description = input.Description
	if err := s.repo.UpdateTermAndConditions(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

/*
Mengambil semua TAC.
Mengembalikan list TAC atau error jika gagal.
*/
func (s *termsAndConditionsService) GetAllTermAndConditions() ([]*model.TermsAndConditionsModel, error) {
	return s.repo.FindAllTermAndConditions()
}

/*
Mendefinisikan operasi service TAC.
Menyediakan method untuk menambah, update, dan ambil TAC dengan hasil sukses atau error.
*/
type TermsAndConditionsService interface {
	AddTermsAndConditions(input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	GetAllTermAndConditions() ([]*model.TermsAndConditionsModel, error)
	UpdateTermAndConditions(id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	UpdateTermsAndConditionsByUserID(userID string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
}

/*
Membuat service TAC.
Mengembalikan instance TermsAndConditionsService yang siap digunakan.
*/
func NewTermsAndConditionsService(repo repository.TermsAndConditionsRepository) TermsAndConditionsService {
	return &termsAndConditionsService{repo: repo}
}
