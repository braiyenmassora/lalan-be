package route

import (
	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
	"net/http"
)

/*
Mendaftarkan route terms and conditions.
Menyiapkan endpoint untuk menambah TAC dengan autentikasi.
*/
func TermsAndConditionsRoutes(h *handler.TermsAndConditionsHandler) {
	v1 := "/v1/tnc"

	addHandler := middleware.JWTMiddleware(http.HandlerFunc(h.AddTermsAndConditions))
	updateHandler := middleware.JWTMiddleware(http.HandlerFunc(h.UpdateTermAndConditions))

	// TermsAndConditions protected routes
	http.Handle(v1+"/add", addHandler)
	http.Handle(v1+"/update", updateHandler)
}
