package message

const (
	// Pesan sukses
	MsgHosterCreatedSuccess = "Hoster created successfully."   // Pesan sukses pembuatan hoster
	MsgHosterLoginSuccess   = "Hoster logged in successfully." // Pesan sukses login hoster

	// Pesan error
	MsgHosterCreatedFailed      = "Failed to create hoster."   // Pesan gagal pembuatan hoster
	MsgHosterLoginFailed        = "Hoster login failed."       // Pesan gagal login hoster
	MsgHosterInvalidEmail       = "Invalid email format."      // Pesan format email tidak valid
	MsgHosterInvalidCredentials = "Invalid email or password." // Pesan kredensial tidak valid

	// Pesan validasi
	MsgHosterEmailExists  = "A hoster with this email already exists." // Pesan email hoster sudah ada
	MsgHosterWeakPassword = "The provided password is too weak."       // Pesan password terlalu lemah
	MsgHosterNotFound     = "Hoster not found."                        // Pesan hoster tidak ditemukan

	// Pesan umum
	MsgInternalServerError = "Internal server error." // Pesan error server internal
	MsgBadRequest          = "Bad request."           // Pesan request tidak valid
	MsgUnauthorized        = "Unauthorized access."   // Pesan akses tidak diizinkan
	MsgNotAllowed          = "Action not allowed."    // Pesan aksi tidak diizinkan

	// Pesan kategori
	MsgCategoryNameExists     = "Category name already exists." // Pesan nama kategori sudah ada
	MsgCategoryNotFound       = "Category not found."           // Pesan kategori tidak ditemukan
	MsgCategoryCreatedSuccess = "Category created successfully" // Pesan sukses pembuatan kategori

	// Pesan Item
	MsgItemNameExists     = "Item name already exists." // Pesan nama item sudah ada
	MsgItemNotFound       = "Item not found."           // Pesan item tidak ditemukan
	MsgItemCreatedSuccess = "Item created successfully" // Pesan sukses pembuatan item
)
