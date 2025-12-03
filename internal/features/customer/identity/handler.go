package identity

import (
	"log"
	"net/http"
	"strings"

	"lalan-be/internal/message"
	"lalan-be/internal/middleware"
	"lalan-be/internal/response"
)

/*
IdentityHandler adalah HTTP transport layer untuk fitur verifikasi identitas (KTP).
Tanggung jawab handler TERBATAS pada:
• Validasi method, header, dan multipart form
• Ekstrak userID dari context (middleware)
• Validasi format file dasar
• Memanggil IdentityService
• Mapping error ke HTTP status + response
Seluruh logika bisnis dan penyimpanan tetap berada di service/repository.
*/
type IdentityHandler struct {
	service *IdentityService
}

/*
NewIdentityHandler membuat instance handler dengan dependency service yang sudah disuntikkan.
Digunakan saat setup router.
*/
func NewIdentityHandler(service *IdentityService) *IdentityHandler {
	return &IdentityHandler{service: service}
}

/*
UploadKTP menangani endpoint POST /identity/ktp (upload KTP).

Catatan: setiap upload diperlakukan sebagai entri baru — kita tidak akan
menghapus atau memperbarui record lama. Ini menjaga referensi historis
dan menjamin foreign key yang menunjuk ke record lama tetap valid.

Alur kerja:
1. Validasi method POST
2. Ambil userID dari context (middleware JWT)
3. Parse multipart form (max 10 MB)
4. Validasi field "ktp" ada dan bertipe image/*
5. Panggil service.UploadKTP() untuk proses upload + simpan ke storage + DB

Output sukses:
- Status: 200 OK
- Body:   null
- Message: "KTP berhasil diupload"

Output error:
- 400 Bad Request  → method salah / form tidak valid / bukan gambar
- 401 Unauthorized → token tidak valid / userID kosong
- 500 Internal      → kegagalan storage / database
*/
func (h *IdentityHandler) UploadKTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Identity.UploadKTP: received request")

	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	// Ambil userID dari middleware context agar konsisten dengan handler lain
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	// Parse multipart form (max 10 MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Identity.UploadKTP: failed to parse multipart form: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	file, header, err := r.FormFile("ktp")
	if err != nil {
		log.Printf("Identity.UploadKTP: missing ktp file: %v", err)
		response.BadRequest(w, message.KTPRequired)
		return
	}
	defer file.Close()

	// Validasi tipe file harus gambar
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		response.BadRequest(w, "file must be an image")
		return
	}

	// Delegasi ke service (service hanya boleh simpan untuk userID yang sama)
	if err := h.service.UploadKTP(r.Context(), userID, file); err != nil {
		log.Printf("Identity.UploadKTP: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, nil, message.KTPUploaded)
}

/*
UpdateKTP menangani endpoint PUT /identity/ktp (re-upload KTP oleh customer).

Catatan: re-upload = entri baru. Handler ini mengunggah file baru dan service
akan menyimpan record baru di tabel `identity`. Record lama tidak diubah
atau dihapus.

Alur kerja:
1. Validasi method PUT
2. Ambil userID dari context
3. Pastikan Content-Type multipart/form-data
4. Parse form dan ambil file "ktp"
5. Panggil service.UpdateKTP() → file baru disimpan, status identity di-reset ke "pending"

Output sukses:
- Status: 200 OK
- Body:   null
- Message: "KTP berhasil diupdate dan status disetel menjadi 'pending'"

Output error:
- 400 Bad Request  → method salah / header salah / file tidak ada
- 401 Unauthorized → token tidak valid
- 500 Internal      → error storage / database
*/
func (h *IdentityHandler) UpdateKTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Identity.UpdateKTP: received request")

	if r.Method != http.MethodPut {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	// Ambil userID dari middleware context agar konsisten dengan handler lain
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
		response.BadRequest(w, "Content-Type must be multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Identity.UpdateKTP: failed to parse multipart form: %v", err)
		response.BadRequest(w, message.BadRequest)
		return
	}

	file, _, err := r.FormFile("ktp")
	if err != nil {
		log.Printf("Identity.UpdateKTP: missing ktp file: %v", err)
		response.BadRequest(w, message.KTPRequired)
		return
	}
	defer file.Close()

	if err := h.service.UpdateKTP(r.Context(), userID, file); err != nil {
		log.Printf("Identity.UpdateKTP: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	response.OK(w, nil, message.KTPUpdated)
}

/*
GetStatusKTP menangani endpoint GET /identity/ktp/status
Mengembalikan status verifikasi KTP milik user yang sedang login.

Alur kerja:
1. Validasi method GET
2. Ambil userID dari context
3. Panggil service.GetStatusKTP() untuk mengambil record identity terakhir

Output sukses:
- Status: 200 OK
- Body:   object status KTP (id, ktp_url, status, reason, dll)
- Message: "Status KTP berhasil diambil"

Output error:
- 401 Unauthorized → token tidak valid
- 404 Not Found     → user belum pernah upload KTP
- 500 Internal      → error service / repository
*/
func (h *IdentityHandler) GetStatusKTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Identity.GetStatusKTP: received request")

	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, message.MethodNotAllowed)
		return
	}

	// Ambil userID dari middleware context agar konsisten dengan handler lain
	userID := middleware.GetUserID(r)
	if userID == "" {
		response.Unauthorized(w, message.Unauthorized)
		return
	}

	status, err := h.service.GetStatusKTP(r.Context(), userID)
	if err != nil {
		log.Printf("Identity.GetStatusKTP: service error: %v", err)
		response.Error(w, http.StatusInternalServerError, message.InternalError)
		return
	}

	if status == nil {
		response.Error(w, http.StatusNotFound, message.NotFound)
		return
	}

	response.OK(w, status, message.KTPStatusRetrieved)
}
