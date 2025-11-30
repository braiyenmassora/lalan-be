package message

const (
	// General messages
	Success          = "success"
	BadRequest       = "invalid request"
	Unauthorized     = "unauthorized"
	Forbidden        = "forbidden"
	MethodNotAllowed = "method not allowed"
	InternalError    = "internal error"

	// Entity-specific messages
	Created       = "%s created"
	Updated       = "%s updated"
	Deleted       = "%s deleted"
	NotFound      = "%s not found"
	AlreadyExists = "%s exists"
	Required      = "%s required"
	InvalidFormat = "invalid %s format"
	TooLong       = "%s too long"
	FileTooLarge  = "%s too large"
	UploadFailed  = "failed to upload %s"

	// ITEM specific messages (optional, explicit)
	ItemCreated  = "item created"
	ItemUpdated  = "item updated"
	ItemDeleted  = "item deleted"
	ItemNotFound = "item not found"

	// Authentication & Authorization
	LoginFailed            = "invalid email or password"
	AdminAccessRequired    = "admin access required"
	HosterAccessRequired   = "hoster access required"
	CustomerAccessRequired = "customer access required"
	InvalidStatus          = "invalid status"
	BookingConflict        = "booking conflict"

	// OTP RELATED
	OTPSent             = "OTP sent to %s"
	OTPResent           = "OTP resent to %s"
	OTPExpired          = "OTP expired"
	OTPInvalid          = "invalid OTP"
	OTPAttemptsExceeded = "too many attempts, request new OTP"
	OTPAlreadyVerified  = "email already verified"

	// BOOKING RELATED
	BookingCreated          = "booking has been created"
	BookingConfirmed        = "booking confirmed"
	BookingCancelled        = "booking cancelled"
	BookingAlreadyConfirmed = "booking already confirmed"
	BookingAlreadyCancelled = "booking already cancelled"
	BookingNotCancellable   = "booking cannot be cancelled"
	BookingOverlap          = "booking time overlaps"

	// IDENTITY VERIFICATION
	IdentitySubmitted               = "identity submitted"
	IdentityApproved                = "identity approved"
	IdentityRejected                = "identity rejected: %s"
	IdentityAlreadyVerified         = "identity already verified"
	IdentityAlreadyUploaded         = "identity already uploaded"
	IdentityCanOnlyUpdateIfRejected = "can only update if rejected"
	IdentityPendingReview           = "identity under review"

	// Additional
	EmailAlreadyExists    = "email already exists"
	CategoryAlreadyExists = "category already exists"
)
