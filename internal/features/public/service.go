package public

import (
	"errors"

	"lalan-be/internal/model"
	"lalan-be/pkg/message"
)

/*
publicService menyediakan logika bisnis untuk data publik.
Menggunakan repository untuk akses data tanpa autentikasi.
*/
type publicService struct {
	repo PublicRepository
}

/*
Methods untuk publicService menangani operasi bisnis kategori, item, dan terms publik.
Dipanggil oleh handler untuk endpoint umum.
*/
func (s *publicService) GetAllCategory() ([]*model.CategoryModel, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return categories, nil
}

func (s *publicService) GetAllItems() ([]*model.ItemModel, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return items, nil
}

func (s *publicService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	tacs, err := s.repo.GetAllTermsAndConditions()
	if err != nil {
		return nil, errors.New(message.MsgInternalServerError)
	}
	return tacs, nil
}

/*
PublicService mendefinisikan kontrak operasi bisnis publik.
Diimplementasikan oleh publicService.
*/
type PublicService interface {
	GetAllCategory() ([]*model.CategoryModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
}

/*
NewPublicService membuat instance PublicService.
Menginisialisasi service dengan repository.
*/
func NewPublicService(repo PublicRepository) PublicService {
	return &publicService{repo: repo}
}
