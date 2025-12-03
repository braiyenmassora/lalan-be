package public

import (
	"errors"

	"lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
publicService mengatur logika bisnis untuk fitur publik.
Berperan sebagai lapisan tengah antara handler dan repository.
*/
type publicService struct {
	repo PublicRepository
}

/*
GetAllCategory mengambil semua kategori publik.

Langkah:
1. Panggil repository untuk ambil data model
2. Mapping dari model ke DTO

Output:
- ([]dto.CategoryDTO, nil) jika sukses
- (nil, error) jika terjadi kesalahan internal
*/
func (s *publicService) GetAllCategory() ([]dto.CategoryPublicResponse, error) {
	categories, err := s.repo.GetAllCategory()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	// Mapping Model -> DTO
	dtos := make([]dto.CategoryPublicResponse, 0)
	for _, cat := range categories {
		dtos = append(dtos, dto.CategoryPublicResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		})
	}

	return dtos, nil
}

/*
GetAllItems mengambil semua item publik.

Langkah:
1. Panggil repository untuk ambil data model
2. Mapping dari model ke DTO

Output:
- ([]dto.ItemDTO, nil) jika sukses
- (nil, error) jika terjadi kesalahan internal
*/
func (s *publicService) GetAllItems() ([]dto.ItemPublicResponse, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	// Mapping Model -> DTO
	dtos := make([]dto.ItemPublicResponse, 0)
	for _, item := range items {
		dtos = append(dtos, dto.ItemPublicResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Photos:      item.Photos,
			Stock:       item.Stock,
			PickupType:  string(item.PickupType),
			PricePerDay: item.PricePerDay,
			Deposit:     item.Deposit,
			Discount:    item.Discount,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			CategoryID:  item.CategoryID,
			HosterID:    item.HosterID,
		})
	}

	return dtos, nil
}

/*
GetAllTermsAndConditions mengambil semua syarat & ketentuan publik.

Langkah:
1. Panggil repository untuk ambil data model
2. Mapping dari model ke DTO

Output:
- ([]dto.TermsAndConditionsDTO, nil) jika sukses
- (nil, error) jika terjadi kesalahan internal
*/
func (s *publicService) GetAllTermsAndConditions() ([]dto.TermsAndConditionsPublicResponse, error) {
	tacs, err := s.repo.GetAllTermsAndConditions()
	if err != nil {
		return nil, errors.New(message.InternalError)
	}

	// Mapping Model -> DTO
	dtos := make([]dto.TermsAndConditionsPublicResponse, 0)
	for _, tac := range tacs {
		dtos = append(dtos, dto.TermsAndConditionsPublicResponse{
			ID:          tac.ID,
			Description: tac.Description,
			CreatedAt:   tac.CreatedAt,
			UpdatedAt:   tac.UpdatedAt,
			UserID:      tac.UserID,
		})
	}

	return dtos, nil
}

/*
GetItemDetail mengambil detail lengkap item dengan JOIN.

Parameter:
- itemID: UUID item yang ingin diambil

Langkah:
1. Panggil repository untuk ambil data JOIN (sudah dalam format DTO)
2. Return langsung karena repository sudah mapping ke DTO

Output:
- (*dto.ItemDetailResponse, nil) jika sukses
- (nil, message.ItemNotFound) jika item tidak ditemukan
- (nil, message.InternalError) jika terjadi kesalahan internal
*/
func (s *publicService) GetItemDetail(itemID string) (*dto.ItemDetailResponse, error) {
	itemDetail, err := s.repo.GetItemDetail(itemID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(message.ItemNotFound)
		}
		return nil, errors.New(message.InternalError)
	}

	return itemDetail, nil
}

/*
PublicService adalah kontrak untuk logika bisnis fitur publik.
Digunakan oleh handler untuk dependency injection.
*/
type PublicService interface {
	GetAllCategory() ([]dto.CategoryPublicResponse, error)
	GetAllItems() ([]dto.ItemPublicResponse, error)
	GetAllTermsAndConditions() ([]dto.TermsAndConditionsPublicResponse, error)
	GetItemDetail(itemID string) (*dto.ItemDetailResponse, error)
}

/*
NewPublicService membuat instance service dengan menyuntikkan repository.

Output:
- PublicService siap digunakan
*/
func NewPublicService(repo PublicRepository) PublicService {
	return &publicService{repo: repo}
}
