package repository

import (
	"fmt"
	"time"

	"github.com/touchsung/spd-fiber-booking-system/models"
)

type MockRepository struct {
	defaultBookings map[string]*models.Booking
}

func NewMockRepository() *MockRepository {
	mockRepo := &MockRepository{
		defaultBookings: make(map[string]*models.Booking),
	}

	baseTime := time.Now().Add(-24 * time.Hour) // Start from yesterday
	// Initialize default bookings (ID 1-10)
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("%d", i)
		mockRepo.defaultBookings[id] = &models.Booking{
			ID:        id,
			UserID:    fmt.Sprintf("user%d", i),
			ServiceID: fmt.Sprintf("service%d", i),
			Price:     float64(i * 10000), // Some will be high-value
			Status:    models.StatusConfirmed,
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour), // Spread over time
		}
	}

	return mockRepo
}

func (m *MockRepository) GetBooking(bookingID string) (*models.Booking, bool) {
	booking, exists := m.defaultBookings[bookingID]
	return booking, exists
}

func (m *MockRepository) GetAllBookings() []*models.Booking {
	bookings := make([]*models.Booking, 0, len(m.defaultBookings))
	for _, booking := range m.defaultBookings {
		bookings = append(bookings, booking)
	}
	return bookings
}

func (m *MockRepository) UpdateBookingStatus(bookingID string, status models.BookingStatus) bool {
	if booking, exists := m.defaultBookings[bookingID]; exists {
		booking.Status = status
		return true
	}
	return false
}

func (m *MockRepository) ClearBookings() {
	m.defaultBookings = make(map[string]*models.Booking)
}

func (m *MockRepository) SaveBooking(booking *models.Booking) {
	m.defaultBookings[booking.ID] = booking
}
