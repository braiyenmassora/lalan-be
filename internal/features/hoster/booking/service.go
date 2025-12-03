package booking

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	dto "lalan-be/internal/dto"
	"lalan-be/internal/message"
)

/*
BookingService adalah kontrak untuk logika bisnis domain booking dari perspektif hoster.
Hanya menyediakan operasi read (list & detail) — hoster tidak bisa membuat booking.
*/
type BookingService interface {
	GetListBookings(hosterID string) ([]dto.BookingListByHosterResponse, error)
	// GetCustomerList returns customers who placed bookings on the given hoster.
	GetCustomerList(hosterID string) ([]dto.CustomerListByHosterResponse, error)
	GetDetailBooking(hosterID, bookingID string) (*dto.BookingDetailByHosterResponse, error)
	UpdateBookingStatus(hosterID, bookingID, newStatus string) error
}

/*
bookingService adalah implementasi konkret dari BookingService.
Mengandung dependency ke repository untuk akses data.
*/
type bookingService struct {
	repo HosterBookingRepository
}

/*
NewBookingService membuat instance service dengan dependency injection.

Output:
- BookingService siap digunakan
*/
func NewBookingService(repo HosterBookingRepository) BookingService {
	return &bookingService{repo: repo}
}

/*
GetListBookings mengambil daftar ringkas semua booking milik hoster.

Alur kerja:
1. Validasi hosterID tidak kosong
2. Panggil repository
3. Wrap error menjadi InternalError

Output sukses:
- ([]dto.BookingListHosterDTO, nil)
Output error:
- (nil, error) → unauthorized / internal error
*/
func (s *bookingService) GetListBookings(hosterID string) ([]dto.BookingListByHosterResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	bookings, err := s.repo.GetListBookings(hosterID)
	if err != nil {
		log.Printf("GetListBookings(hoster service): repo error for hoster %s: %v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return bookings, nil
}

// GetCustomerList returns the list of customers who have made bookings on hosterID.
func (s *bookingService) GetCustomerList(hosterID string) ([]dto.CustomerListByHosterResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	list, err := s.repo.GetCustomerList(hosterID)
	if err != nil {
		log.Printf("GetCustomerList(hoster service): repo error hoster=%s err=%v", hosterID, err)
		return nil, errors.New(message.InternalError)
	}

	return list, nil
}

/*
GetDetailBooking mengambil detail lengkap satu booking milik hoster.

Alur kerja:
1. Validasi hosterID tidak kosong
2. Ambil detail dari repository
3. Cek apakah booking ditemukan
4. Authorization: pastikan hosterID sesuai dengan Booking.HosterID

Output sukses:
- (*dto.BookingDetailHosterDTO, nil)
Output error:
- (nil, error) → unauthorized / not found / internal error
*/
func (s *bookingService) GetDetailBooking(hosterID, bookingID string) (*dto.BookingDetailByHosterResponse, error) {
	if hosterID == "" {
		return nil, errors.New(message.Unauthorized)
	}

	detail, err := s.repo.GetBookingDetail(bookingID)
	if err != nil {
		// Normalize DB not-found to a friendly NotFound error for handlers.
		if err == sql.ErrNoRows {
			log.Printf("GetDetailBooking(hoster service): booking not found booking=%s hoster=%s", bookingID, hosterID)
			return nil, fmt.Errorf(message.NotFound, "booking")
		}
		log.Printf("GetDetailBooking(hoster service): repo error for booking %s: %v", bookingID, err)
		return nil, errors.New(message.InternalError)
	}

	if detail == nil {
		return nil, fmt.Errorf(message.NotFound, "booking")
	}

	if detail.Booking.HosterID != hosterID {
		return nil, errors.New(message.Unauthorized)
	}

	return detail, nil
}

/*
UpdateBookingStatus mengupdate status booking dengan validasi business logic.

Parameter:
- hosterID: UUID hoster yang login (untuk authorization)
- bookingID: UUID booking yang akan diupdate
- newStatus: Status baru (on_progress, on_rent, completed)

Alur kerja:
1. Validasi hosterID tidak kosong
2. Ambil status booking saat ini
3. Validasi authorization: pastikan booking milik hoster ini
4. Validasi status transition (sequential only)
5. Update status di repository

Valid transitions:
- pending → on_progress
- on_progress → on_rent
- on_rent → completed

Output sukses:
- nil
Output error:
- error → unauthorized / not found / invalid transition / internal error
*/
func (s *bookingService) UpdateBookingStatus(hosterID, bookingID, newStatus string) error {
	if hosterID == "" {
		return errors.New(message.Unauthorized)
	}

	// Validasi newStatus harus salah satu dari status yang valid
	validStatuses := map[string]bool{
		"on_progress": true,
		"on_rent":     true,
		"completed":   true,
	}
	if !validStatuses[newStatus] {
		return errors.New(message.InvalidStatus)
	}

	// Ambil detail booking untuk authorization dan validasi
	detail, err := s.repo.GetBookingDetail(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("UpdateBookingStatus: booking not found: %s", bookingID)
			return fmt.Errorf(message.NotFound, "booking")
		}
		log.Printf("UpdateBookingStatus: repo error: %v", err)
		return errors.New(message.InternalError)
	}

	// Authorization: pastikan booking milik hoster ini
	if detail.Booking.HosterID != hosterID {
		log.Printf("UpdateBookingStatus: unauthorized hoster=%s booking_hoster=%s", hosterID, detail.Booking.HosterID)
		return errors.New(message.Unauthorized)
	}

	currentStatus := detail.Booking.Status

	// Validasi status transition (must be sequential)
	validTransitions := map[string]string{
		"pending":     "on_progress",
		"on_progress": "on_rent",
		"on_rent":     "completed",
	}

	expectedNext, exists := validTransitions[currentStatus]
	if !exists {
		log.Printf("UpdateBookingStatus: invalid current status: %s", currentStatus)
		return errors.New(message.InvalidStatus)
	}

	if newStatus != expectedNext {
		log.Printf("UpdateBookingStatus: invalid transition from %s to %s (expected %s)", currentStatus, newStatus, expectedNext)
		return errors.New(message.InvalidStatus)
	}

	// Update status di repository
	err = s.repo.UpdateBookingStatus(bookingID, newStatus)
	if err != nil {
		log.Printf("UpdateBookingStatus: update error: %v", err)
		return errors.New(message.InternalError)
	}

	log.Printf("UpdateBookingStatus: success booking=%s %s→%s", bookingID, currentStatus, newStatus)
	return nil
}
