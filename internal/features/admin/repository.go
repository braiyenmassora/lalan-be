package admin

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Implementasi repository admin dengan koneksi database.
*/
type adminRepository struct {
	db *sqlx.DB
}

/*
Membuat admin baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *adminRepository) CreateAdmin(admin *model.AdminModel) error {
	query := `
        INSERT INTO admin (
            id, full_name, email, password_hash
        )
        VALUES (:id, :full_name, :email, :password_hash)
    `
	_, err := r.db.NamedExec(query, admin)
	if err != nil {
		log.Printf("Error inserting admin: %v", err)
		return err // Return the error instead of nil
	}
	return nil
}

/*
Mencari admin berdasarkan email untuk login.
Mengembalikan data admin atau nil jika tidak ditemukan.
*/
func (r *adminRepository) FindByEmailAdminForLogin(email string) (*model.AdminModel, error) {
	var admin model.AdminModel
	query := `
        SELECT
            id, email, password_hash, full_name,
            created_at, updated_at
        FROM admin  
        WHERE email = $1
          AND password_hash IS NOT NULL
        LIMIT 1
    `

	err := r.db.Get(&admin, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {

		return nil, err
	}

	return &admin, nil
}

/*
Mengambil admin berdasarkan ID.
Mengembalikan data admin atau nil jika tidak ditemukan.
*/
func (r *adminRepository) GetAdminByID(id string) (*model.AdminModel, error) {
	var admin model.AdminModel
	query := `
        SELECT
            id, full_name, email, created_at, updated_at
        FROM admin  
        WHERE id = $1
        LIMIT 1
    `

	err := r.db.Get(&admin, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

/*
Membuat kategori baru di database.
Mengembalikan error jika penyisipan gagal.
*/
func (r *adminRepository) CreateCategory(category *model.CategoryModel) error {
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
Memperbarui kategori.
Mengembalikan error jika update gagal atau tidak ada baris terpengaruh.
*/
func (r *adminRepository) UpdateCategory(category *model.CategoryModel) error {
	query := `UPDATE category 
	          SET name = $2, description = $3, updated_at = NOW() 
	          WHERE id = $1`
	result, err := r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("UpdateCategory error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("UpdateCategory RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("UpdateCategory: no rows affected for id %s", category.ID)
	}

	return nil
}

/*
Menghapus kategori berdasarkan ID.
Mengembalikan error jika penghapusan gagal atau tidak ada baris terpengaruh.
*/
func (r *adminRepository) DeleteCategory(id string) error {
	query := "DELETE FROM category WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteCategory error: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DeleteCategory RowsAffected error: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("DeleteCategory: no rows affected for id %s", id)
	}

	return nil
}

/*
Mengambil semua kategori.
Mengembalikan daftar kategori atau error jika gagal.
*/
func (r *adminRepository) GetListCategory() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("GetListCategory error: %v", err)
		return nil, err
	}
	return categories, nil
}

/*
Mendefinisikan operasi repository untuk admin.
Menyediakan method untuk membuat, mencari, dan mengambil admin dengan hasil sukses atau error.
*/
type AdminRepository interface {
	CreateAdmin(admin *model.AdminModel) error
	FindByEmailAdminForLogin(email string) (*model.AdminModel, error)
	GetAdminByID(id string) (*model.AdminModel, error)
	CreateCategory(category *model.CategoryModel) error
	UpdateCategory(category *model.CategoryModel) error
	DeleteCategory(id string) error
	GetListCategory() ([]*model.CategoryModel, error)
}

/*
Membuat repository admin.
Mengembalikan instance AdminRepository yang siap digunakan.
*/
func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}
