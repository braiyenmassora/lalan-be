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

	// ITEM
	ItemCreated           = "item created"
	ItemUpdated           = "item updated"
	ItemDeleted           = "item deleted"
	ItemNotFound          = "item not found"
	ItemRetrieved         = "item retrieved successfully"
	ItemHasActiveBookings = "item masih terikat dengan pesanan aktif dan tidak dapat dihapus"
	ItemVisibilityUpdated = "item visibility updated successfully"

	// Authentication & Authorization
	LoginFailed            = "invalid email or password"
	AdminAccessRequired    = "admin access required"
	HosterAccessRequired   = "hoster access required"
	CustomerAccessRequired = "customer access required"
	InvalidStatus          = "invalid status"
	BookingConflict        = "booking conflict"
	AdminCreated           = "admin created successfully"
	HosterCreated          = "hoster created successfully"
	CustomerCreated        = "customer created successfully"

	// OTP
	OTPSent             = "OTP sent to %s"
	OTPResent           = "OTP resent to %s"
	OTPExpired          = "OTP expired"
	OTPInvalid          = "invalid OTP"
	OTPAttemptsExceeded = "too many attempts, request new OTP"
	OTPAlreadyVerified  = "email already verified"

	// BOOKING
	BookingCreated          = "booking has been created"
	BookingConfirmed        = "booking confirmed"
	BookingCancelled        = "booking cancelled"
	BookingAlreadyConfirmed = "booking already confirmed"
	BookingAlreadyCancelled = "booking already cancelled"
	BookingNotCancellable   = "booking cannot be cancelled"
	BookingOverlap          = "booking time overlaps"
	BookingStatusUpdated    = "booking status updated successfully"

	// KTP
	KTPUploaded                = "KTP uploaded successfully"
	KTPUpdated                 = "KTP updated successfully"
	KTPStatusRetrieved         = "KTP status retrieved"
	KTPListRetrieved           = "KTP list retrieved successfully"
	KTPSubmitted               = "KTP submitted"
	KTPApproved                = "KTP approved"
	KTPRejected                = "KTP rejected"
	KTPRejectedUploadNew       = "KTP rejected, please upload new KTP"
	KTPAlreadyVerified         = "KTP already verified"
	KTPAlreadyUploaded         = "KTP already uploaded"
	KTPCanOnlyUpdateIfRejected = "can only update if rejected"
	KTPPendingReview           = "KTP under review"
	KTPRequired                = "KTP upload required"
	KTPUploadFailed            = "failed to upload KTP"

	// ID Validation
	UserIDRequired     = "user ID required"
	HosterIDRequired   = "hoster ID required"
	CustomerIDRequired = "customer ID required"
	AdminIDRequired    = "admin ID required"
	FileRequired       = "file required"

	// Additional
	EmailAlreadyExists    = "email already exists"
	EmailNotVerified      = "email not verified"
	CategoryAlreadyExists = "category already exists"
	CustomerNotFound      = "customer not found"
	HosterNotFound        = "hoster not found"
	ResetTokenSent        = "reset password token sent"
	ResetTokenInvalid     = "invalid or expired reset token"
	PasswordResetSuccess  = "password reset successfully"

	// PROFILE
	ProfileRetrieved = "profile retrieved successfully"
	ProfileUpdated   = "profile updated successfully"
	ProfileNotFound  = "profile not found"

	// TERMS AND CONDITIONS (TnC)
	TnCCreated   = "terms and conditions created successfully"
	TnCUpdated   = "terms and conditions updated successfully"
	TnCRetrieved = "terms and conditions retrieved successfully"
	TnCNotFound  = "terms and conditions not found"
)
