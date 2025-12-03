// ===================================================================
// File: booking.go
// Deskripsi: Entity Booking - semua yang terkait pemesanan/sewa item
// Catatan: SEMUA model booking HANYA di file ini. Customer & Hoster pakai model yang sama!
// ===================================================================

package domain

import "time"

// ===================================================================
// BOOKING (Header/Master)
// ===================================================================

// Booking adalah entity utama untuk transaksi pemesanan item.
// Satu booking bisa berisi banyak item yang disewa dalam periode yang sama.
//
// Flow status:
// - "pending": Baru dibuat, menunggu pembayaran (locked selama 30 menit)
// - "confirmed": Sudah bayar, menunggu approval hoster
// - "approved": Hoster setuju, booking aktif
// - "rejected": Hoster tolak
// - "completed": Sewa selesai, item sudah dikembalikan
// - "cancelled": Dibatalkan oleh customer/hoster
//
// Relasi:
// - Booking belongs to Customer (user_id)
// - Booking belongs to Hoster (hoster_id) → hoster pemilik item pertama di booking
// - Booking has many BookingItem
// - Booking has one BookingCustomer (snapshot data customer saat booking)
// - Booking may have one Identity (KTP yang dipakai untuk verifikasi)
type Booking struct {
	ID                   string    `json:"id" db:"id"`
	HosterID             string    `json:"hoster_id" db:"hoster_id"`         // Hoster pemilik item (diambil dari item pertama)
	LockedUntil          time.Time `json:"locked_until" db:"locked_until"`   // Waktu kadaluarsa pembayaran (30 menit dari create)
	TimeRemainingMinutes int       `json:"time_remaining_minutes" db:"-"`    // Sisa waktu dalam menit (dihitung runtime, tidak disimpan)
	StartDate            time.Time `json:"start_date" db:"start_date"`       // Tanggal mulai sewa
	EndDate              time.Time `json:"end_date" db:"end_date"`           // Tanggal selesai sewa
	TotalDays            int       `json:"total_days" db:"total_days"`       // Durasi sewa dalam hari
	DeliveryType         string    `json:"delivery_type" db:"delivery_type"` // "pickup" atau "delivery"
	Rental               int       `json:"rental" db:"rental"`               // Total biaya sewa (sum dari semua item)
	Deposit              int       `json:"deposit" db:"deposit"`             // Total deposit (sum dari semua item)
	Discount             int       `json:"discount" db:"discount"`           // Diskon (jika ada)
	Total                int       `json:"total" db:"total"`                 // rental + deposit - discount
	Outstanding          int       `json:"outstanding" db:"outstanding"`     // Sisa yang harus dibayar (awalnya sama dengan total)
	UserID               string    `json:"user_id" db:"user_id"`             // ID customer yang booking
	IdentityID           *string   `json:"identity_id" db:"identity_id"`     // ID KTP yang dipakai (nullable, diisi jika sudah verified)
	Status               string    `json:"status" db:"status"`               // Status booking (lihat keterangan di atas)
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// ===================================================================
// BOOKING ITEM (Detail Item dalam Booking)
// ===================================================================

// BookingItem adalah entity untuk detail item yang disewa dalam satu booking.
// Satu booking bisa punya banyak BookingItem.
//
// Kenapa ada field Name, PricePerDay, dll padahal bisa join ke Item table?
// → Snapshot! Jika harga item berubah di masa depan, booking lama tetap pakai harga lama.
//
// Relasi:
// - BookingItem belongs to Booking
// - BookingItem references Item (item_id) tapi simpan snapshot data untuk history
type BookingItem struct {
	ID              string   `json:"id" db:"id"`
	BookingID       string   `json:"booking_id" db:"booking_id"`
	ItemID          string   `json:"item_id" db:"item_id"`                   // ID item asli (untuk tracking)
	Name            string   `json:"name" db:"name"`                         // Snapshot nama item saat booking dibuat
	Description     string   `json:"description" db:"description"`           // Enriched dari item table (bukan snapshot)
	Photos          []string `json:"photos" db:"photos"`                     // Enriched dari item table (bukan snapshot)
	Quantity        int      `json:"quantity" db:"quantity"`                 // Jumlah unit yang disewa
	PricePerDay     int      `json:"price_per_day" db:"price_per_day"`       // Snapshot harga per hari
	DepositPerUnit  int      `json:"deposit_per_unit" db:"deposit_per_unit"` // Snapshot deposit per unit
	SubtotalRental  int      `json:"subtotal_rental" db:"subtotal_rental"`   // quantity × price_per_day × total_days
	SubtotalDeposit int      `json:"subtotal_deposit" db:"subtotal_deposit"` // quantity × deposit_per_unit
}

// ===================================================================
// BOOKING CUSTOMER (Snapshot Data Customer)
// ===================================================================

// BookingCustomer adalah snapshot data customer pada saat booking dibuat.
// Disimpan terpisah agar perubahan profil customer tidak mengubah data booking lama.
//
// Contoh kasus:
// - Customer booking tgl 1 Jan dengan alamat "Jakarta"
// - Tgl 5 Jan customer update profil, alamat jadi "Bandung"
// - Booking tanggal 1 Jan tetap tampil alamat "Jakarta" (dari snapshot)
//
// Relasi:
// - BookingCustomer belongs to Booking (one-to-one)
type BookingCustomer struct {
	ID        string `json:"id" db:"id"`
	BookingID string `json:"booking_id" db:"booking_id"`
	Name      string `json:"name" db:"name"`       // Nama penerima (bisa beda dengan nama customer di profil)
	Phone     string `json:"phone" db:"phone"`     // No HP penerima
	Email     string `json:"email" db:"email"`     // Email penerima
	Address   string `json:"address" db:"address"` // Alamat pengiriman
	Notes     string `json:"notes" db:"notes"`     // Catatan tambahan dari customer
}
