package repository

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

// CategoryRepository interface
type CategoryRepository interface {
	FindCategoryName(name string) (*model.CategoryModel, error) // Mencari kategori berdasarkan nama
	CreateCategory(category *model.CategoryModel) error         // Membuat kategori baru
}

// categoryRepository struct
type categoryRepository struct {
	db *sqlx.DB // Koneksi database
}

// NewCategoryRepository constructor
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindCategoryName method
func (r *categoryRepository) FindCategoryName(name string) (*model.CategoryModel, error) {
	// Cari kategori berdasarkan nama
	query := "SELECT id, name, description FROM categories WHERE name = $1 LIMIT 1"
	row := r.db.QueryRow(query, name)
	var category model.CategoryModel
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	if err == sql.ErrNoRows {
		log.Printf("FindCategoryName: no rows for name %s", name)
		return nil, nil
	}
	if err != nil {
		log.Printf("FindCategoryName: scan error %v", err)
		return nil, err
	}
	log.Printf("FindCategoryName: found category %s", category.Name)
	return &category, nil
}

// CreateCategory method
func (r *categoryRepository) CreateCategory(category *model.CategoryModel) error {
	// Insert data kategori
	query := "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("CreateCategory: insert error %v", err)
	}
	return err
}
