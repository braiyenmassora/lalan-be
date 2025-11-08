package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"strings"

	"github.com/google/uuid"
)

// CategoryService adalah interface untuk menangani operasi kategori
type CategoryService interface {
	AddCategory(input *model.CategoryModel) (*model.CategoryModel, error)               // Menambahkan kategori baru
	GetAllCategories() ([]*model.CategoryModel, error)                                  // Mendapatkan semua kategori
	GetCategoryByID(id string) (*model.CategoryModel, error)                            // Mendapatkan kategori berdasarkan ID
	UpdateCategory(id string, input *model.CategoryModel) (*model.CategoryModel, error) // Mengupdate kategori
	DeleteCategory(id string) error                                                     // Menghapus kategori
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
func (s *categoryService) AddCategory(input *model.CategoryModel) (*model.CategoryModel, error) {
	// Validasi input
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Memeriksa apakah nama kategori sudah ada
	existing, err := s.repo.FindCategoryName(input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgCategoryNameExists)
	}

	// Menghasilkan ID unik
	input.ID = uuid.New().String()

	// Menyimpan kategori ke database
	if err := s.repo.CreateCategory(input); err != nil {
		return nil, err
	}

	// Mendapatkan kategori yang baru dibuat untuk mengembalikan data lengkap
	return s.repo.FindByID(input.ID)
}

// GetAllCategories mendapatkan semua kategori
func (s *categoryService) GetAllCategories() ([]*model.CategoryModel, error) {
	return s.repo.FindAll()
}

// GetCategoryByID mendapatkan kategori berdasarkan ID
func (s *categoryService) GetCategoryByID(id string) (*model.CategoryModel, error) {
	if id == "" {
		return nil, errors.New("category ID is required")
	}

	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New(message.MsgCategoryNotFound)
	}

	return category, nil
}

// UpdateCategory mengupdate kategori
func (s *categoryService) UpdateCategory(id string, input *model.CategoryModel) (*model.CategoryModel, error) {
	// Validasi ID
	if id == "" {
		return nil, errors.New("category ID is required")
	}

	// Validasi input
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Cek apakah kategori ada
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgCategoryNotFound)
	}

	// Cek apakah nama baru sudah digunakan kategori lain
	if input.Name != existing.Name {
		duplicate, err := s.repo.FindCategoryName(input.Name)
		if err != nil {
			return nil, err
		}
		if duplicate != nil {
			return nil, errors.New(message.MsgCategoryNameExists)
		}
	}

	// Update kategori
	input.ID = id
	if err := s.repo.Update(input); err != nil {
		return nil, err
	}

	// Mendapatkan kategori yang sudah diupdate
	return s.repo.FindByID(id)
}

// DeleteCategory menghapus kategori
func (s *categoryService) DeleteCategory(id string) error {
	if id == "" {
		return errors.New("category ID is required")
	}

	// Cek apakah kategori ada
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New(message.MsgCategoryNotFound)
	}

	// Hapus kategori
	return s.repo.Delete(id)
}
