package identity

import (
	"encoding/json"
	"net/http"

	"lalan-be/internal/response"

	"github.com/gorilla/mux"
)

/*
AdminIdentityHandler menangani endpoint admin untuk verifikasi identitas (KTP).
Berfungsi sebagai adapter antara HTTP request dan service layer.
*/
type AdminIdentityHandler struct {
	service *AdminIdentityService
}

/*
NewAdminIdentityHandler membuat instance handler dengan dependency injection.

Output:
- *AdminIdentityHandler siap digunakan
*/
func NewAdminIdentityHandler(service *AdminIdentityService) *AdminIdentityHandler {
	return &AdminIdentityHandler{service: service}
}

/*
GetPendingIdentities menangani GET /api/v1/admin/identities/pending.

Alur kerja:
1. Panggil service untuk ambil semua identitas berstatus pending
2. Return data atau error

Output sukses:
- 200 OK + list identitas pending
Output error:
- 500 Internal Server Error
*/
func (h *AdminIdentityHandler) GetPendingIdentities(w http.ResponseWriter, r *http.Request) {
	identities, err := h.service.GetPendingIdentities()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve pending identities")
		return
	}
	response.OK(w, identities, "Pending identities retrieved successfully")
}

/*
ValidateIdentity menangani POST /api/v1/admin/identities/{userID}/validate.

Alur kerja:
1. Ambil userID dari path parameter
2. Decode request body (status + reason optional)
3. Panggil service untuk approve/reject KTP
4. Return hasil validasi

Output sukses:
- 200 OK + pesan sukses
Output error:
- 400 Bad Request → body tidak valid / status tidak diperbolehkan
*/
func (h *AdminIdentityHandler) ValidateIdentity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Path parameter now is the identity record id (not user id)
	id := vars["id"]

	var req struct {
		Status string `json:"status"`           // "approved" atau "rejected"
		Reason string `json:"reason,omitempty"` // wajib jika rejected
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	err := h.service.ValidateIdentity(id, req.Status, req.Reason)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, nil, "Identity validated successfully")
}

/*
GetIdentity menangani GET /api/v1/admin/identities/{userID}.

Alur kerja:
1. Ambil userID dari path parameter
2. Panggil service untuk detail identitas user tersebut

Output sukses:
- 200 OK + data identitas
Output error:
- 404 Not Found → user tidak punya identitas / belum upload KTP
*/
func (h *AdminIdentityHandler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	identity, err := h.service.GetIdentity(userID)
	if err != nil {
		response.NotFound(w, "Identity not found")
		return
	}

	response.OK(w, identity, "Identity retrieved successfully")
}
