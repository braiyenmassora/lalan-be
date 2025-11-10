package message

/*
Mendefinisikan konstanta pesan untuk respons API.
Menyediakan pesan standar untuk sukses, error, dan validasi berdasarkan domain.
*/
const (
	// General
	MsgSuccess             = "Operation completed successfully."
	MsgBadRequest          = "Bad request."
	MsgUnauthorized        = "Unauthorized access."
	MsgNotAllowed          = "Action not allowed."
	MsgInternalServerError = "Internal server error."
	MsgUserIDRequired      = "User ID is required."

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
	MsgHosterFetched            = "Hoster data retrieved successfully"

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

	MsgItemStockInvalid       = "Item stock cannot be negative."
	MsgItemPricePerDayInvalid = "Item price per day cannot be negative."
	MsgItemDepositInvalid     = "Item deposit cannot be negative."

	// Term and Conditions
	MsgTermAndConditionsCreatedSuccess      = "Term and conditions created successfully."
	MsgTermAndConditionsUpdatedSuccess      = "Term and conditions updated successfully."
	MsgTermAndConditionsDeletedSuccess      = "Term and conditions deleted successfully."
	MsgTermAndConditionsAlreadyExists       = "Term and conditions description already exists."
	MsgTermAndConditionsNotFound            = "Term and conditions not found."
	MsgTermAndConditionsDescriptionRequired = "Term and conditions description is required."
	MsgTermAndConditionsDescriptionTooShort = "Term and conditions description must be 255 characters."
	MsgTermAndConditionsIDRequired          = "Term and conditions ID is required."
	MsgTermAndConditionsStockInvalid        = "Term and conditions stock cannot be negative."
	MsgTermAndConditionsPricePerDayInvalid  = "Term and conditions price per day cannot be negative."
	MsgTermAndConditionsDepositInvalid      = "Term and conditions deposit cannot be negative."
)
