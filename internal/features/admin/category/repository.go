package category

import (
	"errors"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/domain"
	"lalan-be/internal/message"
)

/*
CategoryRepository adalah kontrak untuk database operations kategori.
*/
type CategoryRepository interface {
	CreateCategory(category *domain.Category) error
	GetAllCategory() ([]domain.Category, error)
	GetCategoryByID(id string) (*domain.Category, error)
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id string) error
}

/*
categoryRepository adalah implementasi konkret.
*/
type categoryRepository struct {
	db *sqlx.DB
}

/*
NewCategoryRepository membuat instance repository.
*/
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

/*
CreateCategory menyimpan kategori baru ke database.

Output:
- error jika duplicate name atau insert gagal.
- nil jika berhasil.
*/
func (r *categoryRepository) CreateCategory(category *domain.Category) error {
	query := `
		INSERT INTO category (name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		log.Printf("CreateCategory (admin): %v", err)
		// Check duplicate name constraint
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return errors.New(message.CategoryAlreadyExists)
		}
		return err
	}
	return nil
}

/*
GetAllCategory mengambil semua kategori.

Output:
- Slice kategori dan error
*/
func (r *categoryRepository) GetAllCategory() ([]domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM category
		ORDER BY name ASC
	`
	var categories []domain.Category
	if err := r.db.Select(&categories, query); err != nil {
		log.Printf("GetAllCategory (admin): %v", err)
		return nil, err
	}
	return categories, nil
}

/*
GetCategoryByID mengambil kategori berdasarkan ID.

Output:
- Pointer ke Category atau nil jika tidak ditemukan
*/
func (r *categoryRepository) GetCategoryByID(id string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM category
		WHERE id = $1
	`
	var category domain.Category
	if err := r.db.Get(&category, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Printf("GetCategoryByID (admin): %v", err)
		return nil, err
	}
	return &category, nil
}

/*
UpdateCategory memperbarui kategori.

Output:
- error jika duplicate name atau update gagal.
- nil jika berhasil.
*/
func (r *categoryRepository) UpdateCategory(category *domain.Category) error {
	query := `
		UPDATE category
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	res, err := r.db.Exec(query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		log.Printf("UpdateCategory (admin): %v", err)
		// Check duplicate name constraint
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return errors.New(message.CategoryAlreadyExists)
		}
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("category not found")
	}
	return nil
}

/*
DeleteCategory menghapus kategori berdasarkan ID.

Output:
- error jika delete gagal.
- nil jika berhasil.
*/
func (r *categoryRepository) DeleteCategory(id string) error {
	query := `DELETE FROM category WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteCategory (admin): %v", err)
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("category not found")
	}
	return nil
}
