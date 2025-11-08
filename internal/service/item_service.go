package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"strings"

	"github.com/google/uuid"
)

type ItemService interface {
	AddItem(input *model.ItemModel) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetItemByID(id string) (*model.ItemModel, error)
	GetItemsByUserID(userID string) ([]*model.ItemModel, error)
	UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error)
	DeleteItem(id string, userID string) error
}

type itemService struct {
	repo repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

// AddItem menambahkan item baru
func (s *itemService) AddItem(input *model.ItemModel) (*model.ItemModel, error) {
	// Validasi input
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New("item name is required")
	}
	if input.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if input.CategoryID == "" {
		return nil, errors.New("category ID is required")
	}
	if input.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if input.PricePerDay < 0 {
		return nil, errors.New("price per day cannot be negative")
	}
	if input.Deposit < 0 {
		return nil, errors.New("deposit cannot be negative")
	}

	// Cek apakah nama item sudah ada untuk user ini
	existing, err := s.repo.FindItemNameByUserID(input.Name, input.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgItemNameExists)
	}

	// Generate ID unik
	input.ID = uuid.New().String()

	// Simpan item ke database
	if err := s.repo.CreateItem(input); err != nil {
		return nil, err
	}

	// Return item yang baru dibuat
	return s.repo.FindByID(input.ID)
}

// GetAllItems mendapatkan semua item
func (s *itemService) GetAllItems() ([]*model.ItemModel, error) {
	return s.repo.FindAll()
}

// GetItemByID mendapatkan item berdasarkan ID
func (s *itemService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New("item ID is required")
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

// GetItemsByUserID mendapatkan semua item milik user tertentu
func (s *itemService) GetItemsByUserID(userID string) ([]*model.ItemModel, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.repo.FindByUserID(userID)
}

// UpdateItem mengupdate item
func (s *itemService) UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error) {
	// Validasi ID
	if id == "" {
		return nil, errors.New("item ID is required")
	}
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Validasi input
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New("item name is required")
	}
	if input.CategoryID == "" {
		return nil, errors.New("category ID is required")
	}
	if input.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if input.PricePerDay < 0 {
		return nil, errors.New("price per day cannot be negative")
	}
	if input.Deposit < 0 {
		return nil, errors.New("deposit cannot be negative")
	}

	// Cek apakah item ada dan milik user ini
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own items")
	}

	// Cek apakah nama baru sudah digunakan item lain dari user ini
	if input.Name != existing.Name {
		duplicate, err := s.repo.FindItemNameByUserID(input.Name, userID)
		if err != nil {
			return nil, err
		}
		if duplicate != nil {
			return nil, errors.New(message.MsgItemNameExists)
		}
	}

	// Update item
	input.ID = id
	input.UserID = userID
	if err := s.repo.Update(input); err != nil {
		return nil, err
	}

	// Return item yang sudah diupdate
	return s.repo.FindByID(id)
}

// DeleteItem menghapus item
func (s *itemService) DeleteItem(id string, userID string) error {
	if id == "" {
		return errors.New("item ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Cek apakah item ada dan milik user ini
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return errors.New("unauthorized: you can only delete your own items")
	}

	// Hapus item
	return s.repo.Delete(id)
}
