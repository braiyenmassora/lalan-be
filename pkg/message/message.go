package message

/*
Konstanta pesan untuk respons API.
Dikelompokkan berdasarkan domain: umum, autentikasi, kategori, dan item.
*/
const (
	// General
	MsgSuccess             = "Operation completed successfully."
	MsgBadRequest          = "Bad request."
	MsgUnauthorized        = "Unauthorized access."
	MsgNotAllowed          = "Action not allowed."
	MsgInternalServerError = "Internal server error."

	// Authentication / Hoster
	MsgHosterCreatedSuccess     = "Hoster created successfully."
	MsgHosterCreatedFailed      = "Failed to create hoster."
	MsgHosterLoginSuccess       = "Hoster logged in successfully."
	MsgHosterLoginFailed        = "Hoster login failed."
	MsgHosterInvalidEmail       = "Invalid email format."
	MsgHosterInvalidCredentials = "Invalid email or password."
	MsgHosterEmailExists        = "A hoster with this email already exists."
	MsgHosterWeakPassword       = "The provided password is too weak."
	MsgFailedToHashPassword     = "Failed to hash password."
	MsgHosterNotFound           = "Hoster not found."

	// Category
	MsgCategoryCreatedSuccess = "Category created successfully."
	MsgCategoryUpdatedSuccess = "Category updated successfully."
	MsgCategoryDeletedSuccess = "Category deleted successfully."
	MsgCategoryNameExists     = "Category name already exists."
	MsgCategoryNotFound       = "Category not found."
	MsgCategoryNameRequired   = "Category name is required."
	MsgCategoryNameTooLong    = "Category name must not exceed 255 characters."
	MsgCategoryIDRequired     = "Category ID is required."

	// Item
	MsgItemCreatedSuccess = "Item created successfully."
	MsgItemUpdatedSuccess = "Item updated successfully."
	MsgItemDeletedSuccess = "Item deleted successfully."
	MsgItemNameExists     = "Item name already exists."
	MsgItemNotFound       = "Item not found."
	MsgItemNameRequired   = "Item name is required."
	MsgItemNameTooLong    = "Item name must not exceed 255 characters."
	MsgItemIDRequired     = "Item ID is required."
)
