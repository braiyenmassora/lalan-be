package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"strings"

	"github.com/google/uuid"
)

/*
Mendefinisikan operasi service item.
Menyediakan method untuk menambah, ambil semua, ambil by ID, ambil by user, update, dan hapus item dengan hasil sukses atau error.
*/
type ItemService interface {
	AddItem(userID string, input *model.ItemModel) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetItemByID(id string) (*model.ItemModel, error)
	GetItemsByUserID(userID string) ([]*model.ItemModel, error)
	UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error)
	DeleteItem(id string, userID string) error
}

/*
Implementasi service item dengan repository.
*/
type itemService struct {
	repo repository.ItemRepository
}

/*
Membuat service item.
Mengembalikan instance ItemService yang siap digunakan.
*/
func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

/*
Menambahkan item baru.
Mengembalikan data item yang dibuat atau error jika validasi/gagal.
*/
func (s *itemService) AddItem(userID string, input *model.ItemModel) (*model.ItemModel, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}
	if userID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}
	if input.CategoryID == "" {
		return nil, errors.New(message.MsgCategoryIDRequired)
	}
	if input.Stock < 0 {
		return nil, errors.New(message.MsgItemStockInvalid)
	}
	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgItemPricePerDayInvalid)
	}
	if input.Deposit < 0 {
		return nil, errors.New(message.MsgItemDepositInvalid)
	}

	existing, err := s.repo.FindItemNameByUserID(input.Name, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgItemNameExists)
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateItem(input); err != nil {
		return nil, err
	}

	return s.repo.FindByID(input.ID)
}

/*
Mengambil semua item.
Mengembalikan daftar item atau error jika gagal.
*/
func (s *itemService) GetAllItems() ([]*model.ItemModel, error) {
	return s.repo.FindAll()
}

/*
Mengambil item berdasarkan ID.
Mengembalikan data item atau error jika tidak ditemukan.
*/
func (s *itemService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}

	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}

	return item, nil
}

/*
Mengambil item berdasarkan user ID.
Mengembalikan daftar item user atau error jika gagal.
*/
func (s *itemService) GetItemsByUserID(userID string) ([]*model.ItemModel, error) {
	if userID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}

	return s.repo.FindByUserID(userID)
}

/*
Mengupdate item.
Mengembalikan data item yang diupdate atau error jika validasi/not found/unauthorized/gagal.
*/
func (s *itemService) UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}
	if userID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}
	if input.CategoryID == "" {
		return nil, errors.New(message.MsgCategoryIDRequired)
	}
	if input.Stock < 0 {
		return nil, errors.New(message.MsgItemStockInvalid)
	}
	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgItemPricePerDayInvalid)
	}
	if input.Deposit < 0 {
		return nil, errors.New(message.MsgItemDepositInvalid)
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if input.Name != existing.Name {
		duplicate, err := s.repo.FindItemNameByUserID(input.Name, userID)
		if err != nil {
			return nil, err
		}
		if duplicate != nil {
			return nil, errors.New(message.MsgItemNameExists)
		}
	}

	input.ID = id
	input.UserID = userID
	if err := s.repo.Update(input); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

/*
Menghapus item berdasarkan ID dan user ID.
Mengembalikan error jika validasi/not found/unauthorized/gagal.
*/
func (s *itemService) DeleteItem(id string, userID string) error {
	if id == "" {
		return errors.New(message.MsgItemIDRequired)
	}
	if userID == "" {
		return errors.New(message.MsgUserIDRequired)
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return errors.New(message.MsgUnauthorized)
	}

	return s.repo.Delete(id)
}
