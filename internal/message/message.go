package message

const (
	// General messages
	Success          = "success"
	BadRequest       = "bad request"
	Unauthorized     = "unauthorized"
	Forbidden        = "forbidden"
	MethodNotAllowed = "method not allowed"
	InternalError    = "internal server error"

	// Entity-specific messages

	Created       = "%s created successfully"
	Updated       = "%s updated successfully"
	Deleted       = "%s deleted successfully"
	NotFound      = "%s not found"
	AlreadyExists = "%s already exists"
	Required      = "%s is required"
	InvalidFormat = "%s format is invalid"
	TooLong       = "%s too long"
	FileTooLarge  = "%s file is too large"
	UploadFailed  = "failed to upload %s"

	// Authentication and authorization messages
	LoginFailed                     = "invalid email or password"
	AdminAccessRequired             = "admin access required"
	HosterAccessRequired            = "hoster access required"
	CustomerAccessRequired          = "customer access required"
	IdentityAlreadyUploaded         = "identity already uploaded"
	IdentityCanOnlyUpdateIfRejected = "identity can only be updated if rejected"
	BookingConflict                 = "booking conflict with existing reservation"
	InvalidStatus                   = "invalid status"
)
