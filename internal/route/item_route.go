package route

import (
	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
	"net/http"
)

// ItemRoutes mengatur semua rute untuk item dengan JWT protection
func ItemRoutes(h *handler.ItemHandler) {
	v1 := "/v1/item"

	// Semua endpoint item membutuhkan autentikasi
	// Endpoint yang membutuhkan JWT
	addHandler := middleware.JWTMiddleware(http.HandlerFunc(h.AddItem))
	myItemsHandler := middleware.JWTMiddleware(http.HandlerFunc(h.GetMyItems))
	updateHandler := middleware.JWTMiddleware(http.HandlerFunc(h.UpdateItem))
	deleteHandler := middleware.JWTMiddleware(http.HandlerFunc(h.DeleteItem))

	// Register protected routes
	http.Handle(v1+"/add", addHandler)
	http.Handle(v1+"/my-items", myItemsHandler)
	http.Handle(v1+"/update", updateHandler)
	http.Handle(v1+"/delete", deleteHandler)

	// Public endpoints (tidak perlu login)
	http.HandleFunc(v1+"/list", h.GetAllItems)
	http.HandleFunc(v1+"/detail", h.GetItemByID)
}
