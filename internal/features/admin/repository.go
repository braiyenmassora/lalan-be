package admin

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

/*
type adminRepository struct
menyediakan akses ke operasi database untuk admin
*/
type adminRepository struct {
	db *sqlx.DB
}

/*
CreateAdmin
membuat admin baru dengan data yang diberikan
*/
func (r *adminRepository) CreateAdmin(admin *model.AdminModel) error {
	/*
	  CreateAdmin query
	  membuat admin baru dengan data yang diberikan
	*/
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

/*
FindByEmailAdminForLogin
mencari admin berdasarkan email untuk login
*/
func (r *adminRepository) FindByEmailAdminForLogin(email string) (*model.AdminModel, error) {
	var admin model.AdminModel
	/*
	  FindByEmailAdminForLogin query
	  mencari admin berdasarkan email untuk login
	*/
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

/*
CreateCategory
membuat kategori baru dengan data yang diberikan
*/
func (r *adminRepository) CreateCategory(category *model.CategoryModel) error {
	/*
	  CreateCategory query
	  membuat kategori baru dengan data yang diberikan
	*/
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

/*
UpdateCategory
memperbarui kategori berdasarkan ID
*/
func (r *adminRepository) UpdateCategory(category *model.CategoryModel) error {
	/*
	  UpdateCategory query
	  memperbarui kategori berdasarkan ID
	*/
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

/*
DeleteCategory
menghapus kategori berdasarkan ID
*/
func (r *adminRepository) DeleteCategory(id string) error {
	/*
	  DeleteCategory query
	  menghapus kategori berdasarkan ID
	*/
	query := `DELETE FROM category WHERE id = $1`
	_, err := r.db.Exec(query, id)
	log.Printf("DeleteCategory: deleted category with ID %s", id)
	return err
}

/*
FindCategoryByName
mencari kategori berdasarkan nama
*/
func (r *adminRepository) FindCategoryByName(name string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	/*
	  FindCategoryByName query
	  mencari kategori berdasarkan nama
	*/
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

/*
FindCategoryByNameExceptID
mencari kategori berdasarkan nama kecuali ID tertentu
*/
func (r *adminRepository) FindCategoryByNameExceptID(name string, id string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	/*
	  FindCategoryByNameExceptID query
	  mencari kategori berdasarkan nama kecuali ID tertentu
	*/
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

/*
GetAllCategory
mengambil semua kategori diurutkan berdasarkan waktu pembuatan
*/
func (r *adminRepository) GetAllCategory() ([]*model.CategoryModel, error) {
	/*
	  GetAllCategory query
	  mengambil semua kategori diurutkan berdasarkan waktu pembuatan
	*/
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
UpdateIdentity
memperbarui data identitas berdasarkan ID untuk approval admin
*/
func (r *adminRepository) UpdateIdentity(identity *model.IdentityModel) error {
	/*
	  UpdateIdentity query
	  memperbarui data identitas berdasarkan ID untuk approval admin
	*/
	query := `
		UPDATE identity
		SET
			ktp_url = $1,
			verified = $2,
			status = $3,
			reason = $4,
			verified_at = $5,
			updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Exec(query, identity.KTPURL, identity.Verified, identity.Status, identity.Reason, identity.VerifiedAt, identity.ID)
	if err != nil {
		log.Printf("UpdateIdentity: error updating identity: %v", err)
		return err
	}
	log.Printf("UpdateIdentity: updated identity %s", identity.ID)
	return nil
}

/*
GetIdentityCustomerByID
mengambil data identitas berdasarkan ID
*/
func (r *adminRepository) GetIdentityCustomerByID(identityID string) (*model.IdentityModel, error) {
	var identity model.IdentityModel
	/*
	  GetIdentityCustomerByID query
	  mengambil data identitas berdasarkan ID
	*/
	query := `
		SELECT
			id,
			user_id,
			ktp_url,
			verified,
			status,
			reason,
			verified_at,
			created_at,
			updated_at
		FROM identity
		WHERE id = $1
	`
	err := r.db.Get(&identity, query, identityID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetIdentityCustomerByID: no identity found for ID %s", identityID)
			return nil, nil
		}
		log.Printf("GetIdentityCustomerByID: error querying ID %s: %v", identityID, err)
		return nil, err
	}
	log.Printf("GetIdentityCustomerByID: found identity for ID %s", identityID)
	return &identity, nil
}

/*
UpdateIdentityStatus
memperbarui status identitas berdasarkan ID
*/
func (r *adminRepository) UpdateIdentityStatus(identityID string, status string, rejectedReason string, verified bool, verifiedAt *time.Time) error {
	/*
	  UpdateIdentityStatus query
	  memperbarui status identitas berdasarkan ID
	*/
	query := `
		UPDATE identity
		SET
			status = $1,
			reason = $2,
			verified = $3,
			verified_at = $4,
			updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(query, status, rejectedReason, verified, verifiedAt, identityID)
	if err != nil {
		log.Printf("UpdateIdentityStatus: error updating identity %s: %v", identityID, err)
		return err
	}
	log.Printf("UpdateIdentityStatus: updated identity %s to status %s", identityID, status)
	return nil
}

/*
UpdateBookingIdentityStatusByUserID
memperbarui status identitas di booking berdasarkan user ID
*/
func (r *adminRepository) UpdateBookingIdentityStatusByUserID(userID string, status string) error {
	/*
	  UpdateBookingIdentityStatusByUserID query
	  memperbarui status identitas di booking berdasarkan user ID
	*/
	query := `
		UPDATE booking
		SET
			identity_status = $1,
			updated_at = NOW()
		WHERE user_id = $2
	`
	_, err := r.db.Exec(query, status, userID)
	if err != nil {
		log.Printf("UpdateBookingIdentityStatusByUserID: error updating booking for user %s: %v", userID, err)
		return err
	}
	log.Printf("UpdateBookingIdentityStatusByUserID: updated booking identity status for user %s to %s", userID, status)
	return nil
}

/*
GetIdentityByCustomerID
mengambil data identitas berdasarkan user ID
*/
func (r *adminRepository) GetIdentityByCustomerID(userID string) (*model.IdentityModel, error) {
	var identity model.IdentityModel
	/*
	  GetIdentityByCustomerID query
	  mengambil data identitas berdasarkan user ID
	*/
	query := `
		SELECT
			id,
			user_id,
			ktp_url,
			verified,
			status,
			reason,
			verified_at,
			created_at,
			updated_at
		FROM identity
		WHERE user_id = $1
	`
	err := r.db.Get(&identity, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetIdentityByCustomerID: no identity found for user ID %s", userID)
			return nil, nil
		}
		log.Printf("GetIdentityByCustomerID: error querying user ID %s: %v", userID, err)
		return nil, err
	}
	log.Printf("GetIdentityByCustomerID: found identity for user ID %s", userID)
	return &identity, nil
}

/*
AdminRepository
mendefinisikan kontrak untuk akses data admin dan kategori
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
	UpdateIdentity(identity *model.IdentityModel) error
	GetIdentityCustomerByID(identityID string) (*model.IdentityModel, error)
	UpdateIdentityStatus(identityID string, status string, rejectedReason string, verified bool, verifiedAt *time.Time) error
	UpdateBookingIdentityStatusByUserID(userID string, status string) error
	GetIdentityByCustomerID(userID string) (*model.IdentityModel, error)
}

/*
NewAdminRepository
membuat instance baru AdminRepository dengan database yang diberikan
*/
func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}
