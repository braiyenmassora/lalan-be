package message

// Konstanta pesan untuk respons API standar.
const (
	MsgSuccess                  = "Operation completed successfully."        // Pesan sukses operasi umum
	MsgHosterCreatedSuccess     = "Hoster created successfully."             // Pesan sukses buat hoster
	MsgHosterLoginSuccess       = "Hoster logged in successfully."           // Pesan sukses login hoster
	MsgHosterCreatedFailed      = "Failed to create hoster."                 // Pesan gagal buat hoster
	MsgHosterLoginFailed        = "Hoster login failed."                     // Pesan gagal login hoster
	MsgHosterInvalidEmail       = "Invalid email format."                    // Pesan validasi email tidak valid
	MsgHosterInvalidCredentials = "Invalid email or password."               // Pesan validasi kredensial salah
	MsgHosterEmailExists        = "A hoster with this email already exists." // Pesan validasi email sudah ada
	MsgHosterWeakPassword       = "The provided password is too weak."       // Pesan validasi password lemah
	MsgHosterNotFound           = "Hoster not found."                        // Pesan hoster tidak ditemukan
	MsgInternalServerError      = "Internal server error."                   // Pesan error server internal
	MsgBadRequest               = "Bad request."                             // Pesan request tidak valid
	MsgUnauthorized             = "Unauthorized access."                     // Pesan akses tidak diizinkan
	MsgNotAllowed               = "Action not allowed."                      // Pesan aksi tidak diizinkan
	MsgCategoryNameExists       = "Category name already exists."            // Pesan validasi nama kategori ada
	MsgCategoryNotFound         = "Category not found."                      // Pesan kategori tidak ditemukan
	MsgCategoryCreatedSuccess   = "Category created successfully"            // Pesan sukses buat kategori
	MsgItemNameExists           = "Item name already exists."                // Pesan validasi nama item ada
	MsgItemNotFound             = "Item not found."                          // Pesan item tidak ditemukan
	MsgItemCreatedSuccess       = "Item created successfully"                // Pesan sukses buat item
)
