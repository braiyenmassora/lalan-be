package item

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"lalan-be/internal/config"
	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"
	"lalan-be/internal/utils"
)

/*
ItemService adalah kontrak untuk logika bisnis domain item dari perspektif hoster.
Menyediakan operasi read (list) dan create untuk hoster.
*/
type ItemService interface {
	GetListItem(hosterID string) ([]dto.ItemListByHosterResponse, error)
	GetItemDetail(hosterID, itemID string) (*dto.ItemDetailByHosterResponse, error)
	CreateItem(ctx context.Context, item *domain.Item, photoFiles []*multipart.FileHeader) (*dto.ItemDetailByHosterResponse, error)
	DeleteItem(hosterID, itemID string) error
	UpdateItem(hosterID, itemID string, req *dto.UpdateItemRequestRequest) error // Tambah ini
}

/*
itemService adalah implementasi konkret dari ItemService.
Mengandung dependency ke repository untuk akses data.
*/
type itemService struct {
	repo    HosterItemRepository
	storage utils.Storage
	config  config.StorageConfig // Tambah config untuk akses ItemBucket
}

/*
NewItemService membuat instance service dengan dependency injection.

Output:
- ItemService siap digunakan
*/
func NewItemService(repo HosterItemRepository, storage utils.Storage, config config.StorageConfig) ItemService {
	return &itemService{repo: repo, storage: storage, config: config}
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
GetItemDetail mengambil detail lengkap satu item milik hoster.

Alur kerja:
1. Validasi hosterID dan itemID tidak kosong
2. Panggil repository dengan filter hosterID (memastikan ownership)
3. Wrap error menjadi NotFound atau InternalError

Output sukses:
- (*dto.ItemDetailByHosterResponse, nil)
Output error:
- (nil, error) → unauthorized / not found / internal error
*/
func (s *itemService) GetItemDetail(hosterID, itemID string) (*dto.ItemDetailByHosterResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}
	if itemID == "" {
		return nil, errors.New(message.BadRequest)
	}

	detail, err := s.repo.GetItemDetail(hosterID, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(message.NotFound)
		}
		log.Printf("GetItemDetail(hoster service): repo error for hoster %s item %s: %v", hosterID, itemID, err)
		return nil, errors.New(message.InternalError)
	}

	return detail, nil
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
func (s *itemService) CreateItem(ctx context.Context, item *domain.Item, photoFiles []*multipart.FileHeader) (*dto.ItemDetailByHosterResponse, error) {
	// Validasi existing
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
	// Ubah validasi pickup_type
	log.Printf("CreateItem: item.PickupType='%s', expected self='%s' delivery='%s'", item.PickupType, domain.PickupMethodSelfPickup, domain.PickupMethodDelivery)

	if !(item.PickupType == domain.PickupMethodSelfPickup || item.PickupType == domain.PickupMethodDelivery) {
		log.Printf("CreateItem: invalid pickup_type")
		return nil, errors.New(message.BadRequest)
	}
	// Handle upload jika ada photoFiles
	if len(photoFiles) > 0 {
		var photoURLs []string
		now := time.Now()
		dateStr := now.Format("02122005") // DDMMYYYY, e.g., 02122025
		for i, fileHeader := range photoFiles {
			newFilename := "item" + strconv.Itoa(i+1) + "_" + dateStr // Sederhanakan filename, karena itemID di path
			uploadPath := item.HosterID + "/item/" + item.ID          // Tambah /item/{itemID}
			metadata, err := s.storage.UploadFile(ctx, fileHeader, uploadPath, s.config.HosterBucket, newFilename)
			if err != nil {
				log.Printf("CreateItem: upload failed for %s: %v", fileHeader.Filename, err)
				return nil, errors.New(message.InternalError)
			}
			photoURLs = append(photoURLs, metadata.URL)
		}
		item.Photos = photoURLs
	}

	// Panggil repository
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

	// Ambil photos sebelum delete
	photos, err := s.repo.GetItemPhotos(itemID)
	if err != nil {
		log.Printf("DeleteItem: failed to get photos for item %s: %v", itemID, err)
		// Lanjut delete, atau return error
	}

	// Delete dari DB
	if err := s.repo.DeleteItem(hosterID, itemID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(message.BadRequest)
		}
		log.Printf("DeleteItem(hoster service): repo error hoster=%s item=%s err=%v", hosterID, itemID, err)
		return errors.New(message.InternalError)
	}

	// Hapus photos dari storage
	for _, url := range photos {
		path := utils.ExtractPathFromURL(url, s.config.Domain, s.config.HosterBucket) // Gunakan utils
		if err := s.storage.Delete(context.Background(), path, s.config.HosterBucket); err != nil {
			log.Printf("DeleteItem: failed to delete photo %s: %v", url, err)
		}
	}

	return nil
}

/*
UpdateItem mengubah data item oleh hoster.
Validasi:
  - hosterID tidak kosong -> Unauthorized
  - itemID tidak kosong, req tidak nil -> BadRequest
  - Validasi field req (stock >= 0, pickup_type valid, dll.)

Business:
  - Panggil repo.UpdateItem
*/
func (s *itemService) UpdateItem(hosterID, itemID string, req *dto.UpdateItemRequestRequest) error {
	if hosterID == "" {
		return errors.New(message.Unauthorized)
	}
	if itemID == "" || req == nil {
		return errors.New(message.BadRequest)
	}

	// Validasi field req
	if req.Stock != nil && *req.Stock < 0 {
		return errors.New(message.BadRequest)
	}
	if req.PickupType != nil && !(*req.PickupType == "self_pickup" || *req.PickupType == "delivery") {
		return errors.New(message.BadRequest)
	}
	if req.Deposit != nil && *req.Deposit < 0 {
		return errors.New(message.BadRequest)
	}
	if req.Discount != nil && *req.Discount < 0 {
		return errors.New(message.BadRequest)
	}

	// Panggil repo.UpdateItem
	if err := s.repo.UpdateItem(hosterID, itemID, req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(message.BadRequest)
		}
		log.Printf("UpdateItem(service): repo error hoster=%s item=%s err=%v", hosterID, itemID, err)
		return errors.New(message.InternalError)
	}

	return nil
}
