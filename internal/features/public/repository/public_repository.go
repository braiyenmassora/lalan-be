package repository

import (
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type PublicRepository interface {
	GetListCategory() ([]*model.CategoryModel, error)
}

type publicRepository struct {
	db *sqlx.DB
}

func (r *publicRepository) GetListCategory() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("GetListCategory error: %v", err)
		return nil, err
	}
	return categories, nil
}

func NewPublicRepository(db *sqlx.DB) PublicRepository {
	return &publicRepository{db: db}
}
