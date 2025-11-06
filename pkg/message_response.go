package pkg

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

const (
	// success messages
	MsgHosterCreatedSuccess = "Hoster successfully created."
	MsgHosterLoginSuccess   = "Hoster logged in successfully."

	// error messages
	MsgHosterCreatedFailed      = "Failed to create hoster."
	MsgHosterLoginFailed        = "Hoster login failed."
	MsgHosterInvalidEmail       = "Invalid email format."
	MsgHosterInvalidCredentials = "Invalid email or password."

	// validation messages
	MsgHosterEmailExists  = "Hoster with this email already exists."
	MsgHosterWeakPassword = "The provided password is too weak."
	MsgHosterNotFound     = "Hoster not found."

	// general messages
	MsgInternalServerError = "Internal server error."
	MsgBadRequest          = "Bad request."
	MsgUnauthorized        = "Unauthorized access."

	// category
	MsgCategoryNameExist = "Category name already exist"
)
