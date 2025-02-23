package service

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/touchsung/spd-fiber-booking-system/internal/domain"
	"github.com/touchsung/spd-fiber-booking-system/internal/repository"
	"github.com/touchsung/spd-fiber-booking-system/internal/util"
)

type BookingService struct {
	cache          *repository.InMemoryCache
	mockRepository *repository.MockRepository
}

func NewBookingService(cache *repository.InMemoryCache, mockRepo *repository.MockRepository) *BookingService {
	return &BookingService{
		cache:          cache,
		mockRepository: mockRepo,
	}
}

func (s *BookingService) requiresCreditCheck(price float64) bool {
	return price > 50000
}

func (s *BookingService) generateRandomStatus() domain.BookingStatus {
	if rand.Float64() < 0.5 {
		return domain.StatusRejected
	}
	return domain.StatusConfirmed
}

func (s *BookingService) simulateCreditCheck(bookingID string) domain.CreditCheckResult {
	time.Sleep(2 * time.Second)
	return domain.CreditCheckResult{
		BookingID: bookingID,
		Status:    s.generateRandomStatus(),
	}
}

func (s *BookingService) processCreditCheck(bookingID string) {
	result := s.simulateCreditCheck(bookingID)
	s.cache.UpdateBookingStatus(result.BookingID, result.Status)
}

func checkExpiredTime(date time.Time) bool {
	return date.Before(time.Now().Add(-5 * time.Minute))
}

func (s *BookingService) CreateBooking(request domain.BookingRequest) (*domain.Booking, error) {
	booking := &domain.Booking{
		ID:        util.GenerateID(),
		UserID:    request.UserID,
		ServiceID: request.ServiceID,
		Price:     request.Price,
		Status:    domain.StatusPending,
	}

	s.cache.SaveBooking(booking)

	if s.requiresCreditCheck(booking.Price) {
		go s.processCreditCheck(booking.ID)
	}

	return booking, nil
}

func (s *BookingService) GetBooking(bookingID string) (*domain.Booking, error) {
	// Try to get from cache first
	if booking, exists := s.cache.GetBooking(bookingID); exists {
		return booking, nil
	}

	// If not in cache, try to get from mock repository
	if booking, exists := s.mockRepository.GetBooking(bookingID); exists {
		// Save to cache for future use
		s.cache.SaveBooking(booking)
		return booking, nil
	}

	return nil, fmt.Errorf("booking not found")
}

func (s *BookingService) ListBookings(sortBy *domain.SortOption, highValueOnly *bool) []*domain.Booking {
	// Get all bookings from cache and mock repository
	allBookings := make([]*domain.Booking, 0)
	allBookings = append(allBookings, s.cache.GetAllBookings()...)
	allBookings = append(allBookings, s.mockRepository.GetAllBookings()...)

	// Remove duplicates
	uniqueBookings := make(map[string]*domain.Booking)
	for _, booking := range allBookings {
		uniqueBookings[booking.ID] = booking
	}

	allBookings = make([]*domain.Booking, 0, len(uniqueBookings))
	for _, booking := range uniqueBookings {
		allBookings = append(allBookings, booking)
	}

	// Filter high-value bookings if requested
	if highValueOnly != nil && *highValueOnly {
		filteredBookings := make([]*domain.Booking, 0)
		for _, booking := range allBookings {
			if booking.Price > 50000 {
				filteredBookings = append(filteredBookings, booking)
			}
		}
		allBookings = filteredBookings
	}

	// Sort bookings
	if sortBy != nil {
		switch *sortBy {
		case domain.SortByPrice:
			sort.Slice(allBookings, func(i, j int) bool {
				return allBookings[i].Price < allBookings[j].Price
			})
		case domain.SortByDate:
			sort.Slice(allBookings, func(i, j int) bool {
				return allBookings[i].CreatedAt.Before(allBookings[j].CreatedAt)
			})
		}
	} else {
		// Default sort by ID if sortBy is nil
		sort.Slice(allBookings, func(i, j int) bool {
			id1, err1 := strconv.Atoi(allBookings[i].ID)
			id2, err2 := strconv.Atoi(allBookings[j].ID)
			if err1 != nil || err2 != nil {
				// Fallback to string comparison if conversion fails
				return allBookings[i].ID < allBookings[j].ID
			}
			return id1 < id2
		})
	}

	return allBookings
}

func (s *BookingService) CancelBooking(bookingID string) error {
	// Try to get from cache first
	if booking, exists := s.cache.GetBooking(bookingID); exists {
		if booking.Status == domain.StatusConfirmed {
			return fmt.Errorf("cannot cancel a confirmed booking")
		}
		// Check if booking exists in mock repository before deleting
		if _, repoExists := s.mockRepository.GetBooking(bookingID); !repoExists {
			return fmt.Errorf("booking not found in repository")
		}
		// Delete from cache
		s.cache.DeleteBooking(bookingID)
	}

	// Update status in mock repository
	if exists := s.mockRepository.UpdateBookingStatus(bookingID, domain.StatusCanceled); !exists {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (s *BookingService) CancelExpiredBookings() {
	pendingBookings := s.cache.GetAllBookings()
	for _, booking := range pendingBookings {
		if booking.Status == domain.StatusPending && checkExpiredTime(booking.CreatedAt) {
			s.cache.UpdateBookingStatus(booking.ID, domain.StatusCanceled)
			s.mockRepository.UpdateBookingStatus(booking.ID, domain.StatusCanceled)
		}
	}
}
