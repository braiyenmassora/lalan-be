package public

import (
	"errors"

	"lalan-be/internal/message"
	"lalan-be/internal/model"
)

/*
publicService
mengelola logika bisnis untuk data publik menggunakan repository
*/
type publicService struct {
	repo PublicRepository
}

/*
GetAllCategory
mengambil semua kategori melalui repository
*/
func (s *publicService) GetAllCategory() ([]*model.CategoryModel, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return categories, nil
}

/*
GetAllItems
mengambil semua item melalui repository
*/
func (s *publicService) GetAllItems() ([]*model.ItemModel, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return items, nil
}

/*
GetAllTermsAndConditions
mengambil semua syarat dan ketentuan melalui repository
*/
func (s *publicService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	tacs, err := s.repo.GetAllTermsAndConditions()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}
	return tacs, nil
}

/*
PublicService
interface untuk operasi service publik
*/
type PublicService interface {
	GetAllCategory() ([]*model.CategoryModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
}

/*
NewPublicService
membuat instance PublicService dengan repository
*/
func NewPublicService(repo PublicRepository) PublicService {
	return &publicService{repo: repo}
}
