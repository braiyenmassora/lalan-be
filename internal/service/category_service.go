package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg"

	"github.com/google/uuid"
)

type CategoryService interface {
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// register category baru
func (s *categoryService) AddCategory(input *model.CategoryModel) error {
	// cek category name
	existing, err := s.repo.FindCategoryName(input.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(pkg.MsgCategoryNameExist)
	}
	// generate id
	input.ID = uuid.New().String()
	return s.repo.CreateCategory(input)
}
