package admin

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

/*
adminRepository menyediakan akses ke operasi database untuk admin.
Menggunakan sqlx.DB untuk query dan eksekusi.
*/
type adminRepository struct {
	db *sqlx.DB
}

/*
Methods adminRepository menangani CRUD admin dan kategori.
Menggunakan query SQL untuk interaksi database.
*/
func (r *adminRepository) CreateAdmin(admin *model.AdminModel) error {
	query := `
		INSERT INTO admin (
			email,
			password_hash,
			full_name,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, admin.Email, admin.PasswordHash, admin.FullName, admin.CreatedAt, admin.UpdatedAt).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)
	log.Printf("CreateAdmin: inserted admin with email %s, ID %s", admin.Email, admin.ID)
	return err
}

func (r *adminRepository) FindByEmailAdminForLogin(email string) (*model.AdminModel, error) {
	var admin model.AdminModel
	query := `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			created_at,
			updated_at
		FROM admin
		WHERE email = $1
	`
	err := r.db.Get(&admin, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindByEmailAdminForLogin: no admin found for email %s", email)
			return nil, nil
		}
		log.Printf("FindByEmailAdminForLogin: error querying email %s: %v", email, err)
		return nil, err
	}
	log.Printf("FindByEmailAdminForLogin: found admin for email %s", email)
	return &admin, nil
}

func (r *adminRepository) CreateCategory(category *model.CategoryModel) error {
	query := `
		INSERT INTO category (
			name,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	log.Printf("CreateCategory: inserted category with name %s, ID %s", category.Name, category.ID)
	return err
}

func (r *adminRepository) UpdateCategory(category *model.CategoryModel) error {
	query := `
		UPDATE category
		SET
			name = $1,
			description = $2,
			updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, category.Name, category.Description, category.UpdatedAt, category.ID)
	log.Printf("UpdateCategory: updated category with ID %s", category.ID)
	return err
}

func (r *adminRepository) DeleteCategory(id string) error {
	query := `DELETE FROM category WHERE id = $1`
	_, err := r.db.Exec(query, id)
	log.Printf("DeleteCategory: deleted category with ID %s", id)
	return err
}

func (r *adminRepository) FindCategoryByName(name string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM category
		WHERE name = $1
	`
	err := r.db.Get(&category, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindCategoryByName: no category found for name %s", name)
			return nil, nil
		}
		log.Printf("FindCategoryByName: error querying name %s: %v", name, err)
		return nil, err
	}
	log.Printf("FindCategoryByName: found category for name %s", name)
	return &category, nil
}

func (r *adminRepository) FindCategoryByNameExceptID(name string, id string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM category
		WHERE name = $1 AND id != $2
	`
	err := r.db.Get(&category, query, name, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindCategoryByNameExceptID: no category found for name %s except ID %s", name, id)
			return nil, nil
		}
		log.Printf("FindCategoryByNameExceptID: error querying name %s except ID %s: %v", name, id, err)
		return nil, err
	}
	log.Printf("FindCategoryByNameExceptID: found category for name %s except ID %s", name, id)
	return &category, nil
}

func (r *adminRepository) GetAllCategory() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Printf("GetAllCategory error: %v", err)
		return nil, err
	}
	return categories, nil
}

/*
AdminRepository mendefinisikan kontrak untuk akses data admin dan kategori.
Wajib diimplementasikan oleh semua penyimpanan data.
*/
type AdminRepository interface {
	CreateAdmin(admin *model.AdminModel) error
	FindByEmailAdminForLogin(email string) (*model.AdminModel, error)
	CreateCategory(category *model.CategoryModel) error
	UpdateCategory(category *model.CategoryModel) error
	DeleteCategory(id string) error
	FindCategoryByName(name string) (*model.CategoryModel, error)
	FindCategoryByNameExceptID(name string, id string) (*model.CategoryModel, error)
	GetAllCategory() ([]*model.CategoryModel, error)
}

/*
NewAdminRepository membuat instance baru AdminRepository.
Mengembalikan interface AdminRepository dengan database yang diberikan.
*/
func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}
