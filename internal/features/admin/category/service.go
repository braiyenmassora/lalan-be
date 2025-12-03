package category

import (
	"errors"
	"log"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
CategoryService adalah kontrak untuk business logic kategori.
*/
type CategoryService interface {
	CreateCategory(req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAllCategory() ([]dto.CategoryListResponse, error)
	GetCategoryByID(id string) (*dto.CategoryResponse, error)
	UpdateCategory(id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id string) error
}

/*
categoryService adalah implementasi konkret.
*/
type categoryService struct {
	repo CategoryRepository
}

/*
NewCategoryService membuat instance service.
*/
func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

/*
CreateCategory membuat kategori baru.

Output:
- CategoryResponse jika berhasil
- error jika nama sudah exist atau validasi gagal
*/
func (s *categoryService) CreateCategory(req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Validasi
	if req.Name == "" {
		return nil, errors.New("category name required")
	}

	// Build entity
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to DB
	if err := s.repo.CreateCategory(category); err != nil {
		if err.Error() == message.CategoryAlreadyExists {
			return nil, errors.New(message.CategoryAlreadyExists)
		}
		log.Printf("CreateCategory service: %v", err)
		return nil, errors.New(message.InternalError)
	}

	// Return response
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

/*
GetAllCategory mengambil semua kategori.
*/
func (s *categoryService) GetAllCategory() ([]dto.CategoryListResponse, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		log.Printf("GetAllCategory service: %v", err)
		return nil, errors.New(message.InternalError)
	}

	// Convert to DTO
	var response []dto.CategoryListResponse
	for _, cat := range categories {
		response = append(response, dto.CategoryListResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		})
	}
	return response, nil
}

/*
GetCategoryByID mengambil detail kategori.
*/
func (s *categoryService) GetCategoryByID(id string) (*dto.CategoryResponse, error) {
	if id == "" {
		return nil, errors.New("category ID required")
	}

	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		log.Printf("GetCategoryByID service: %v", err)
		return nil, errors.New(message.InternalError)
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

/*
UpdateCategory memperbarui kategori.
*/
func (s *categoryService) UpdateCategory(id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// Validasi
	if id == "" {
		return nil, errors.New("category ID required")
	}
	if req.Name == "" {
		return nil, errors.New("category name required")
	}

	// Cek apakah kategori ada
	existing, err := s.repo.GetCategoryByID(id)
	if err != nil {
		log.Printf("UpdateCategory service (check): %v", err)
		return nil, errors.New(message.InternalError)
	}
	if existing == nil {
		return nil, errors.New("category not found")
	}

	// Build updated entity
	category := &domain.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	// Update DB
	if err := s.repo.UpdateCategory(category); err != nil {
		if err.Error() == message.CategoryAlreadyExists {
			return nil, errors.New(message.CategoryAlreadyExists)
		}
		if err.Error() == "category not found" {
			return nil, errors.New("category not found")
		}
		log.Printf("UpdateCategory service: %v", err)
		return nil, errors.New(message.InternalError)
	}

	// Ambil data terbaru
	updated, err := s.repo.GetCategoryByID(id)
	if err != nil || updated == nil {
		log.Printf("UpdateCategory service (refetch): %v", err)
		return nil, errors.New(message.InternalError)
	}

	return &dto.CategoryResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

/*
DeleteCategory menghapus kategori.
*/
func (s *categoryService) DeleteCategory(id string) error {
	if id == "" {
		return errors.New("category ID required")
	}

	if err := s.repo.DeleteCategory(id); err != nil {
		if err.Error() == "category not found" {
			return errors.New("category not found")
		}
		log.Printf("DeleteCategory service: %v", err)
		return errors.New(message.InternalError)
	}
	return nil
}
