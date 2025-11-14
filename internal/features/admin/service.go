package admin

import (
	"lalan-be/internal/model"
)

/*
Implementasi service admin dengan repository.
*/
type adminService struct {
	repo AdminRepository
}

/*
Membuat admin baru melalui service.
Mengembalikan error jika penyisipan gagal.
*/
func (s *adminService) CreateAdmin(admin *model.AdminModel) error {
	return s.repo.CreateAdmin(admin)
}

/*
Mencari admin berdasarkan email untuk login melalui service.
Mengembalikan data admin atau nil jika tidak ditemukan.
*/
func (s *adminService) FindByEmailAdminForLogin(email string) (*model.AdminModel, error) {
	return s.repo.FindByEmailAdminForLogin(email)
}

/*
Membuat kategori baru melalui service.
Mengembalikan error jika penyisipan gagal.
*/
func (s *adminService) CreateCategory(category *model.CategoryModel) error {
	return s.repo.CreateCategory(category)
}

/*
Memperbarui kategori melalui service.
Mengembalikan error jika update gagal.
*/
func (s *adminService) UpdateCategory(category *model.CategoryModel) error {
	return s.repo.UpdateCategory(category)
}

/*
Menghapus kategori berdasarkan ID melalui service.
Mengembalikan error jika penghapusan gagal.
*/
func (s *adminService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}

/*
Mendefinisikan operasi service untuk admin.
Menyediakan method untuk mengelola admin dan kategori dengan hasil sukses atau error.
*/
type AdminService interface {
	CreateAdmin(admin *model.AdminModel) error
	FindByEmailAdminForLogin(email string) (*model.AdminModel, error)
	CreateCategory(category *model.CategoryModel) error
	UpdateCategory(category *model.CategoryModel) error
	DeleteCategory(id string) error
}

/*
Membuat service admin.
Mengembalikan instance AdminService yang siap digunakan.
*/
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}
