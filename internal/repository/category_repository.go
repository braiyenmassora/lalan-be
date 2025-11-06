package repository

import (
	"database/sql"
	"lalan-be/internal/model"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(category *model.CategoryModel) error
	FindCategoryName(name string) (*model.CategoryModel, error)
}

// implementasi konkret pakai PostgreSQL (sqlx.DB)
type categoryRepository struct {
	db *sqlx.DB
}

// inisialisasi repository dengan koneksi DB
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// insert data baru ke table category
func (r *categoryRepository) CreateCategory(h *model.CategoryModel) error {
	query := `INSERT INTO category(
		id, name, description
	)
	VALUES(:id :name, :description)
	`
	_, err := r.db.NamedExec(query, h)
	return err
}

func (r *categoryRepository) FindCategoryName(name string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	query := ` SELECT * FROM categories WHERE LOWER(name)= LOWER($1) LIMIT 1`
	err := r.db.Get(&category, query, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}
