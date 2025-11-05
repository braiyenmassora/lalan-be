package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type HosterService interface {
	RegisterHoster(input *model.HosterModel) error
}

type hosterService struct {
	repo repository.HosterRepository
}

func NewHosterService(repo repository.HosterRepository) HosterService {
	return &hosterService{repo: repo}
}

// RegisterHoster membuat akun hoster baru
func (s *hosterService) RegisterHoster(input *model.HosterModel) error {
	// cek email sudah terdaftar atau belum
	existing, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	// generate ID unik
	input.ID = uuid.New().String()

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	input.Password = string(hashed)

	// simpan ke database
	return s.repo.CreateHoster(input)
}
