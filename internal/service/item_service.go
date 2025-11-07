package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"

	"github.com/google/uuid"
)

type ItemService interface {
	AddItem(input *model.ItemModel) error
}

type itemService struct {
	repo repository.ItemRepository
}

func NewItemRepository(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

// Add Item baru
func (s *itemService) AddItem(input *model.ItemModel) error {
	// check item name berdasarkan userID
	existing, err := s.repo.FindItemNameByUserID(input.Name, input.UserID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(message.MsgItemNameExists)

	}
	// menghasilkan id unik
	input.ID = uuid.New().String()
	// simpan data
	return s.repo.CreateItem(input)
}
