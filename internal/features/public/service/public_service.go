package service

import (
	"lalan-be/internal/features/public/repository"
	"lalan-be/internal/model"
)

type PublicService interface {
	GetListCategory() ([]*model.CategoryModel, error)
}

type publicService struct {
	repo repository.PublicRepository
}

func (s *publicService) GetListCategory() ([]*model.CategoryModel, error) {
	return s.repo.GetListCategory()
}

func NewPublicService(repo repository.PublicRepository) PublicService {
	return &publicService{repo: repo}
}
