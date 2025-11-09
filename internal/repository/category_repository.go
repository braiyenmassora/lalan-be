package repository

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
CategoryRepository mendefinisikan operasi untuk mengelola data kategori di database.
*/
type CategoryRepository interface {
	FindCategoryName(name string) (*model.CategoryModel, error)
	CreateCategory(category *model.CategoryModel) error
	FindAll() ([]*model.CategoryModel, error)
	FindByID(id string) (*model.CategoryModel, error)
	Update(category *model.CategoryModel) error
	Delete(id string) error
}

/*
categoryRepository menyimpan koneksi database untuk operasi repository kategori.
*/
type categoryRepository struct {
	db *sqlx.DB
}

/*
NewCategoryRepository membuat instance CategoryRepository dengan koneksi database.
Mengembalikan interface CategoryRepository.
*/
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

/*
FindCategoryName mencari kategori berdasarkan nama di tabel categories.
Mengembalikan pointer ke model dan error; (nil, nil) jika tidak ditemukan.
*/
func (r *categoryRepository) FindCategoryName(name string) (*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE name = $1 LIMIT 1"
	var category model.CategoryModel
	err := r.db.Get(&category, query, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindCategoryName error: %v", err)
		return nil, err
	}
	return &category, nil
}

/*
CreateCategory menyisipkan kategori baru ke tabel categories.
Mengembalikan error jika gagal.
*/
func (r *categoryRepository) CreateCategory(category *model.CategoryModel) error {
	query := `INSERT INTO categories (id, name, description, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("CreateCategory error: %v", err)
		return err
	}
	return nil
}

/*
FindAll mengambil semua kategori dari tabel categories dalam urutan descending.
Mengembalikan slice pointer ke model dan error jika gagal.
*/
func (r *categoryRepository) FindAll() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("FindAll error: %v", err)
		return nil, err
	}
	return categories, nil
}

/*
FindByID mencari kategori berdasarkan ID di tabel categories.
Mengembalikan pointer ke model dan error; (nil, nil) jika tidak ditemukan.
*/
func (r *categoryRepository) FindByID(id string) (*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1 LIMIT 1"
	var category model.CategoryModel
	err := r.db.Get(&category, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
		return nil, err
	}
	return &category, nil
}

/*
Update memperbarui data kategori di tabel categories berdasarkan ID.
Mengembalikan error jika gagal atau tidak ada baris yang terpengaruh.
*/
func (r *categoryRepository) Update(category *model.CategoryModel) error {
	query := `UPDATE categories 
	          SET name = $2, description = $3, updated_at = NOW() 
	          WHERE id = $1`
	result, err := r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("Update error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Update RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Update: no rows affected for id %s", category.ID)
	}

	return nil
}

/*
Delete menghapus kategori dari tabel categories berdasarkan ID.
Mengembalikan error jika gagal atau tidak ada baris yang terpengaruh.
*/
func (r *categoryRepository) Delete(id string) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Delete error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Delete RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Delete: no rows affected for id %s", id)
	}

	return nil
}
