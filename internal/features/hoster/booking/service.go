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
