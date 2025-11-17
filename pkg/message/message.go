package message

/*
Mendefinisikan pesan umum untuk respons API.
Digunakan di seluruh aplikasi untuk error dan sukses standar.
*/
const (
	MsgSuccess             = "Success."
	MsgBadRequest          = "Bad request."
	MsgUnauthorized        = "Unauthorized."
	MsgMethodNotAllowed    = "Method not allowed."
	MsgInternalServerError = "Internal server error."
)

/*
Mendefinisikan pesan autentikasi hoster.
Digunakan dalam handler login dan registrasi hoster.
*/
const (
	MsgHosterCreatedSuccess = "Hoster created successfully."
	MsgHosterLoginSuccess   = "Hoster logged in successfully."
	MsgHosterLoginFailed    = "Invalid email or password."
	MsgHosterEmailExists    = "Email already exists."
	MsgHosterNotFound       = "Hoster not found."
)

/*
Mendefinisikan pesan autentikasi customer.
Digunakan dalam handler login dan registrasi customer.
*/
const (
	MsgCustomerCreatedSuccess = "Customer created successfully."
	MsgCustomerLoginSuccess   = "Customer logged in successfully."
	MsgCustomerLoginFailed    = "Invalid email or password."
	MsgCustomerEmailExists    = "Email already exists."
	MsgCustomerNotFound       = "Customer not found."
)

/*
Mendefinisikan pesan operasi kategori.
Digunakan dalam CRUD kategori.
*/
const (
	MsgCategoryCreatedSuccess = "Category created successfully."
	MsgCategoryUpdatedSuccess = "Category updated successfully."
	MsgCategoryDeletedSuccess = "Category deleted successfully."
	MsgCategoryNotFound       = "Category not found."
	MsgCategoryNameRequired   = "Category name is required."
	MsgCategoryIDRequired     = "Category ID is required."
)

/*
Mendefinisikan pesan operasi item.
Digunakan dalam CRUD item.
*/
const (
	MsgItemCreatedSuccess = "Item created successfully."
	MsgItemUpdatedSuccess = "Item updated successfully."
	MsgItemDeletedSuccess = "Item deleted successfully."
	MsgItemNotFound       = "Item not found."
	MsgItemNameRequired   = "Item name is required."
	MsgItemIDRequired     = "Item ID is required."
)

/*
Mendefinisikan pesan operasi terms and conditions.
Digunakan dalam CRUD TnC.
*/
const (
	MsgTnCCreatedSuccess      = "Terms and conditions created successfully."
	MsgTnCUpdatedSuccess      = "Terms and conditions updated successfully."
	MsgTnCDeletedSuccess      = "Terms and conditions deleted successfully."
	MsgTnCNotFound            = "Terms and conditions not found."
	MsgTnCDescriptionRequired = "Description is required."
	MsgTnCIDRequired          = "Terms and conditions ID is required."
)
