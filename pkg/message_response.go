package hoster

type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	// success messages
	MsgHosterCreatedSuccess = "Hoster successfully created."
	MsgHosterLoginSuccess   = "Hoster logged in successfully."

	// error messages
	MsgHosterCreatedFailed = "Failed to create hoster."
	MsgHosterLoginFailed   = "Hoster login failed."

	// validation messages
	MsgHosterEmailExists  = "Hoster with this email already exists."
	MsgHosterWeakPassword = "The provided password is too weak."
	MsgHosterNotFound     = "Hoster not found."

	// general messages
	MsgInternalServerError = "Internal server error."
	MsgBadRequest          = "Bad request."
	MsgUnauthorized        = "Unauthorized access."
)
