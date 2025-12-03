package booking

import (
	"errors"
	"fmt"
	"log"
	"time"

	"lalan-be/internal/domain"
	"lalan-be/internal/dto"
	"lalan-be/internal/message"

	"github.com/google/uuid"
)

/*
BookingService adalah kontrak (interface) untuk seluruh logika bisnis domain booking.
Layer ini bertanggung jawab atas:
• Ekstraksi user dari context
• Validasi aturan bisnis
• Perhitungan nilai (total, durasi, dll)
• Authorization (kepemilikan data)
• Orkestrasi repository
Tidak boleh ada detail HTTP atau database di sini.
*/
type BookingService interface {
	CreateBooking(userID string, req dto.CreateBookingByCustomerRequest) (*dto.BookingDetailByCustomerResponse, error)
	GetListBookings(userID string) ([]dto.BookingListByCustomerResponse, error)
	GetDetailBooking(userID string, bookingID string) (*dto.BookingDetailByCustomerResponse, error)
}

/*
bookingService adalah implementasi konkret dari BookingService.
Menyimpan dependency ke repository untuk persistensi data.
*/
type bookingService struct {
	repo BookingRepository
}

/*
NewBookingService membuat instance service yang siap digunakan.
Dependency injection dilakukan di sini untuk memudahkan testing.

Output:
- Implementasi BookingService yang terkoneksi ke repository.
*/
func NewBookingService(repo BookingRepository) BookingService {
	return &bookingService{repo: repo}
}

/*
CreateBooking menangani seluruh proses bisnis pembuatan booking baru.

Alur kerja:
1. Ekstrak user ID dari context (via middleware auth)
2. Validasi KTP user sudah ter-upload (via repository)
3. Parse dan hitung durasi sewa (totalDays)
4. Hitung total rental + deposit - discount
5. Generate booking ID dan locked_until (30 menit)
6. Tentukan hoster_id dari item pertama
7. Bangun entity BookingModel, BookingItem[], dan BookingCustomer
8. Persist semua data via repository dalam satu transaksi

Output sukses:
- *dto.BookingDetailByCustomerResponse (detail lengkap booking yang baru dibuat)
Output error:
- message.Unauthorized → 401 (token invalid/missing)
- "silakan upload ktp terlebih dahulu" → 400 (dari repository)
- "hoster tidak dapat ditentukan..." → 400
- Semua error lain → 500 (internal)
*/
func (s *bookingService) CreateBooking(userID string, req dto.CreateBookingByCustomerRequest) (*dto.BookingDetailByCustomerResponse, error) {
	// 1. Validasi userID yang diberikan oleh handler (middleware)
	if userID == "" {
		return nil, errors.New(message.UserIDRequired)
	}

	// 2. Validasi KTP (hanya cek keberadaan, detail validasi di repo)
	identity, err := s.repo.GetIdentityByUserID(userID)
	if err != nil {
		log.Printf("CreateBooking service: repo error checking identity user %s: %v", userID, err)
		return nil, errors.New(message.InternalError)
	}
	if identity == nil {
		// Jika tidak ada identity verifikasi yang valid, beri tahu user untuk
		// mengunggah dan menyelesaikan proses verifikasi KTP terlebih dahulu.
		return nil, errors.New(message.KTPRequired)
	}

	// 2a. Validasi status KTP - jika rejected, tampilkan reason dari admin
	if identity.Status == "rejected" {
		log.Printf("CreateBooking service: user %s has rejected KTP, cannot create booking", userID)
		// Format: "KTP Rejected - {reason dari admin}"
		if identity.Reason != "" {
			return nil, fmt.Errorf("KTP Rejected - %s", identity.Reason)
		}
		return nil, errors.New(message.KTPRejectedUploadNew)
	}

	// 3. Parse tanggal & hitung durasi
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	totalDays := int(endDate.Sub(startDate).Hours() / 24)

	// 4. Hitung total biaya
	var rentalTotal, depositTotal int
	for _, item := range req.Items {
		rentalTotal += item.SubtotalRental
		depositTotal += item.SubtotalDeposit
	}
	total := rentalTotal + depositTotal - req.Discount
	outstanding := total

	// 5. Generate ID dan waktu lock
	bookingID := uuid.New().String()
	lockedUntil := time.Now().Add(30 * time.Minute)

	// 6. Bangun Booking header
	booking := &domain.Booking{
		ID:                   bookingID,
		HosterID:             "", // akan diisi nanti
		LockedUntil:          lockedUntil,
		TimeRemainingMinutes: 30, // akan dihitung ulang di repo
		StartDate:            startDate,
		EndDate:              endDate,
		TotalDays:            totalDays,
		DeliveryType:         req.DeliveryType,
		Rental:               rentalTotal,
		Deposit:              depositTotal,
		Discount:             req.Discount,
		Total:                total,
		Outstanding:          outstanding,
		UserID:               userID,
		IdentityID:           &identity.ID,
		Status:               "pending",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// 7. Tentukan hoster dari item pertama
	if len(req.Items) > 0 {
		hosterID, err := s.repo.GetHosterIDByItemID(req.Items[0].ItemID)
		if err != nil {
			log.Printf("CreateBooking service: failed resolve hoster for item %s: %v", req.Items[0].ItemID, err)
			return nil, errors.New(message.InternalError)
		}
		if hosterID == "" {
			return nil, errors.New(message.HosterIDRequired)
		}
		booking.HosterID = hosterID
	} else {
		return nil, errors.New("at least one item required")
	}

	// 8. Bangun booking items
	items := make([]domain.BookingItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = domain.BookingItem{
			ID:              uuid.New().String(),
			BookingID:       bookingID,
			ItemID:          it.ItemID,
			Name:            it.Name,
			Quantity:        it.Quantity,
			PricePerDay:     it.PricePerDay,
			DepositPerUnit:  it.DepositPerUnit,
			SubtotalRental:  it.SubtotalRental,
			SubtotalDeposit: it.SubtotalDeposit,
		}
	}

	// 9. Bangun customer data
	customer := domain.BookingCustomer{
		ID:        uuid.New().String(),
		BookingID: bookingID,
		Name:      req.Customer.Name,
		Phone:     req.Customer.Phone,
		Email:     req.Customer.Email,
		Address:   req.Customer.Address,
		Notes:     req.Customer.Notes,
	}
	if customer.Address == "" {
		customer.Address = "N/A"
	}

	// 10. Persist via repository
	detail, err := s.repo.CreateBooking(booking, items, customer)
	if err != nil {
		return nil, err // error sudah sesuai konteks (KTP, DB, dll)
	}

	return detail, nil
}

/*
GetBookingsByUserID mengembalikan daftar ringkas semua booking milik user yang login.

Alur kerja:
1. Ekstrak user ID dari context
2. Panggil repository untuk ambil data

Output sukses:
- []dto.BookingListByCustomerResponse (bisa kosong)
Output error:
- message.Unauthorized → 401
- Semua error repo → 500
*/
func (s *bookingService) GetListBookings(userID string) ([]dto.BookingListByCustomerResponse, error) {
	if userID == "" {
		return nil, errors.New(message.UserIDRequired)
	}

	bookings, err := s.repo.GetListBookings(userID)
	if err != nil {
		log.Printf("GetListBookings service: repo error for user %s: %v", userID, err)
		return nil, errors.New(message.InternalError)
	}

	return bookings, nil
}

/*
GetDetailBooking mengembalikan detail lengkap satu booking dengan pengecekan kepemilikan.

Alur kerja:
1. Ekstrak user ID dari context
2. Ambil detail booking via repository
3. Validasi booking benar-benar ada
4. Validasi user adalah pemilik booking (authorization)

Output sukses:
- *dto.BookingDetailByCustomerResponse (lengkap dengan items, customer, KTP status)
Output error:
- message.Unauthorized → 401 (token invalid atau bukan pemilik)
- message.NotFound + "booking" → 404
- Error repository → 500
*/
func (s *bookingService) GetDetailBooking(userID string, bookingID string) (*dto.BookingDetailByCustomerResponse, error) {
	if userID == "" {
		return nil, errors.New(message.UserIDRequired)
	}

	detail, err := s.repo.GetBookingDetail(bookingID)
	if err != nil {
		log.Printf("GetDetailBooking service: repo error for booking %s: %v", bookingID, err)
		return nil, errors.New(message.InternalError)
	}
	if detail == nil {
		return nil, fmt.Errorf(message.NotFound, "booking")
	}

	// Authorization: user harus pemilik booking
	if detail.Booking.UserID != userID {
		return nil, errors.New(message.Unauthorized)
	}

	return detail, nil
}

// DTO Request untuk booking sudah dipindah ke: internal/dto/booking_dto.go
// - dto.CreateBookingByCustomerRequest
// - dto.CreateBookingItemByCustomerRequest
// - dto.CreateBookingCustomerByCustomerRequest
