package repository

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Implementasi repository kategori dengan koneksi database.
*/
type categoryRepository struct {
	db *sqlx.DB
}

/*
Mencari kategori berdasarkan nama.
Mengembalikan data kategori atau nil jika tidak ditemukan.
*/
func (r *categoryRepository) FindCategoryName(name string) (*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category WHERE name = $1 LIMIT 1"
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
Membuat kategori baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *categoryRepository) CreateCategory(category *model.CategoryModel) error {
	query := `INSERT INTO category (id, name, description, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("CreateCategory error: %v", err)
		return err
	}
	return nil
}

/*
Mengambil semua kategori.
Mengembalikan daftar kategori atau error jika gagal.
*/
func (r *categoryRepository) FindAll() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("FindAll error: %v", err)
		return nil, err
	}
	return categories, nil
}

/*
Mencari kategori berdasarkan ID.
Mengembalikan data kategori atau nil jika tidak ditemukan.
*/
func (r *categoryRepository) FindByID(id string) (*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category WHERE id = $1 LIMIT 1"
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
Memperbarui kategori.
Mengembalikan error jika update gagal atau tidak ada baris terpengaruh.
*/
func (r *categoryRepository) Update(category *model.CategoryModel) error {
	query := `UPDATE category 
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
Menghapus kategori berdasarkan ID.
Mengembalikan error jika penghapusan gagal atau tidak ada baris terpengaruh.
*/
func (r *categoryRepository) Delete(id string) error {
	query := "DELETE FROM category WHERE id = $1"
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

/*
Mendefinisikan operasi repository untuk kategori.
Menyediakan method untuk mencari, membuat, update, dan hapus kategori dengan hasil sukses atau error.
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
Membuat repository kategori.
Mengembalikan instance CategoryRepository yang siap digunakan.
*/
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
