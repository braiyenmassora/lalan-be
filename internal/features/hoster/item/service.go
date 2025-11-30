package item

import (
	"database/sql"
	"errors"
	"log"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
ItemService adalah kontrak untuk logika bisnis domain item dari perspektif hoster.
Menyediakan operasi read (list) dan create untuk hoster.
*/
type ItemService interface {
	GetListItem(hosterID string) ([]dto.ItemListByHosterResponse, error)
	CreateItem(item *domain.Item) (*dto.ItemDetailByHosterResponse, error)
	// DeleteItem menghapus item milik hoster berdasarkan id item
	DeleteItem(hosterID, itemID string) error
}

/*
itemService adalah implementasi konkret dari ItemService.
Mengandung dependency ke repository untuk akses data.
*/
type itemService struct {
	repo HosterItemRepository
}

/*
NewItemService membuat instance service dengan dependency injection.

Output:
- ItemService siap digunakan
*/
func NewItemService(repo HosterItemRepository) ItemService {
	return &itemService{repo: repo}
}

/*
GetListItem mengambil daftar ringkas semua item milik hoster.

Alur kerja:
1. Validasi hosterID tidak kosong
2. Panggil repository
3. Wrap error menjadi InternalError

Output sukses:
- ([]dto.ItemListByHosterResponse, nil)
Output error:
- (nil, error) → unauthorized / internal error
*/
func (s *itemService) GetListItem(hosterID string) ([]dto.ItemListByHosterResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	items, err := s.repo.GetListItem(hosterID)
	if err != nil {
		log.Printf("GetListItem(hoster service): repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return items, nil
}

/*
CreateItem menyimpan item baru (oleh hoster).

Alur kerja:
1. Validasi input minimal (hoster/user id, name, stock, price_per_day, pickup_type)
2. Panggil repository untuk menyimpan item
3. Wrap error menjadi InternalError / BadRequest

Output sukses:
- (*dto.ItemDetailByHosterResponse, nil)
Output error:
- (nil, error) → unauthorized / bad request / internal error
*/
func (s *itemService) CreateItem(item *domain.Item) (*dto.ItemDetailByHosterResponse, error) {
	if item == nil || item.HosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}
	if item.Name == "" {
		return nil, errors.New(message.BadRequest)
	}
	if item.Stock < 0 {
		return nil, errors.New(message.BadRequest)
	}
	if item.PricePerDay <= 0 {
		return nil, errors.New(message.BadRequest)
	}
	if !(item.PickupType == domain.PickupMethodSelfPickup || item.PickupType == domain.PickupMethodDelivery) {
		return nil, errors.New(message.BadRequest)
	}

	created, err := s.repo.CreateItem(item)
	if err != nil {
		log.Printf("CreateItem(hoster service): repo error for hoster %s: %v", item.HosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return created, nil
}

/*
DeleteItem menghapus item oleh hoster.
Validasi:
  - hosterID (diambil dari session/token) tidak boleh kosong -> Unauthorized
  - itemID tidak boleh kosong -> BadRequest

Business:
  - panggil repo.DeleteItem (DELETE WHERE id = $1 AND hoster_id = $2)
  - jika sql.ErrNoRows -> BadRequest (item tidak ada atau bukan milik hoster)
  - jika error lain -> InternalError
*/
func (s *itemService) DeleteItem(hosterID, itemID string) error {
	if hosterID == "" {
		return errors.New(message.Unauthorized)
	}
	if itemID == "" {
		return errors.New(message.BadRequest)
	}

	if err := s.repo.DeleteItem(hosterID, itemID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// item tidak ada atau bukan milik hoster
			return errors.New(message.BadRequest)
		}
		log.Printf("DeleteItem(hoster service): repo error hoster=%s item=%s err=%v", hosterID, itemID, err)
		return errors.New(message.InternalError)
	}

	return nil
}
