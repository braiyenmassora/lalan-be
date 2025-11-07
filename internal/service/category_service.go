package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"

	"github.com/google/uuid"
)

// CategoryService adalah interface untuk menangani operasi kategori
type CategoryService interface {
	AddCategory(input *model.CategoryModel) error // Menambahkan kategori baru
}

// categoryService adalah struct yang mengimplementasikan CategoryService
type categoryService struct {
	repo repository.CategoryRepository // Repository untuk operasi kategori
}

// NewCategoryService membuat instance CategoryService
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// AddCategory menambahkan kategori baru
func (s *categoryService) AddCategory(input *model.CategoryModel) error {
	// Memeriksa apakah nama kategori sudah ada
	existing, err := s.repo.FindCategoryName(input.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(message.MsgCategoryNameExists)
	}

	// Menghasilkan ID unik
	input.ID = uuid.New().String()

	// Menyimpan kategori ke database
	return s.repo.CreateCategory(input)
}
