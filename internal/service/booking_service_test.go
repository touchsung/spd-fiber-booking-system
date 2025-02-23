package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/touchsung/spd-fiber-booking-system/internal/domain"
	"github.com/touchsung/spd-fiber-booking-system/internal/repository"
)

func setupTestService() *BookingService {
	cache := repository.NewInMemoryCache()
	mockRepo := repository.NewMockRepository()
	service := NewBookingService(cache, mockRepo)
	service.mockRepository = mockRepo
	service.mockRepository.ClearBookings()
	return service
}

func TestCheckExpiredTime(t *testing.T) {
	// Test with a time that is not expired
	nonExpiredTime := time.Now().Add(-4 * time.Minute)
	assert.False(t, checkExpiredTime(nonExpiredTime), "Expected non-expired time to return false")

	// Test with a time that is expired
	expiredTime := time.Now().Add(-6 * time.Minute)
	assert.True(t, checkExpiredTime(expiredTime), "Expected expired time to return true")
}

func TestRequiresCreditCheck(t *testing.T) {
	service := setupTestService()

	// Test with a price below the threshold
	assert.False(t, service.requiresCreditCheck(40000), "Expected no credit check for price below threshold")

	// Test with a price at the threshold
	assert.False(t, service.requiresCreditCheck(50000), "Expected no credit check for price at threshold")

	// Test with a price above the threshold
	assert.True(t, service.requiresCreditCheck(60000), "Expected credit check for price above threshold")
}

func TestGenerateRandomStatus(t *testing.T) {
	service := setupTestService()

	// Run multiple times to ensure both statuses are generated
	statusCounts := make(map[domain.BookingStatus]int)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		status := service.generateRandomStatus()
		statusCounts[status]++
	}

	// Check that both statuses were generated
	assert.True(t, statusCounts[domain.StatusRejected] > 0, "Expected StatusRejected to be generated")
	assert.True(t, statusCounts[domain.StatusConfirmed] > 0, "Expected StatusConfirmed to be generated")
}

func TestCreateBooking(t *testing.T) {
	service := setupTestService()

	tests := []struct {
		name    string
		request domain.BookingRequest
	}{
		{
			name: "normal booking",
			request: domain.BookingRequest{
				UserID:    "user1",
				ServiceID: "service1",
				Price:     40000,
			},
		},
		{
			name: "high-value booking",
			request: domain.BookingRequest{
				UserID:    "user2",
				ServiceID: "service2",
				Price:     60000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booking, err := service.CreateBooking(tt.request)

			assert.NoError(t, err)
			assert.NotEmpty(t, booking.ID)
			assert.Equal(t, tt.request.UserID, booking.UserID)
			assert.Equal(t, tt.request.ServiceID, booking.ServiceID)
			assert.Equal(t, tt.request.Price, booking.Price)
			assert.Equal(t, domain.StatusPending, booking.Status)
		})
	}
}

func TestGetBooking(t *testing.T) {
	service := setupTestService()

	// Create a test booking
	request := domain.BookingRequest{
		UserID:    "user1",
		ServiceID: "service1",
		Price:     40000,
	}
	booking, _ := service.CreateBooking(request)

	// Test cases
	t.Run("existing booking", func(t *testing.T) {
		found, err := service.GetBooking(booking.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, booking.ID, found.ID)
	})

	t.Run("non-existent booking", func(t *testing.T) {
		found, err := service.GetBooking("non-existent-id")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestListBookings(t *testing.T) {
	// Setup mock data
	service := setupTestService()

	// Add bookings to cache
	service.cache.SaveBooking(&domain.Booking{
		ID:        "1",
		UserID:    "user1",
		ServiceID: "service1",
		Price:     60000,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		Status:    domain.StatusPending,
	})

	// Add bookings to mock repository
	service.mockRepository.SaveBooking(&domain.Booking{
		ID:        "2",
		UserID:    "user2",
		ServiceID: "service2",
		Price:     40000,
		CreatedAt: time.Now().Add(-48 * time.Hour),
		Status:    domain.StatusPending,
	})

	// Test without filters
	bookings := service.ListBookings(nil, nil)
	assert.Equal(t, 2, len(bookings), "Expected 2 bookings")
	assert.Equal(t, "1", bookings[0].ID, "Expected booking ID 1")

	// Test high value filter
	highValueOnly := true
	bookings = service.ListBookings(nil, &highValueOnly)
	assert.Equal(t, 1, len(bookings), "Expected 1 high-value booking")
	assert.Equal(t, "1", bookings[0].ID, "Expected booking ID 1")

	// Test sorting by price
	sortByPrice := domain.SortByPrice
	bookings = service.ListBookings(&sortByPrice, nil)
	assert.Equal(t, "2", bookings[0].ID, "Expected booking ID 2 to be first when sorted by price")

	// Test sorting by date
	sortByDate := domain.SortByDate
	bookings = service.ListBookings(&sortByDate, nil)
	assert.Equal(t, "2", bookings[0].ID, "Expected booking ID 2 to be first when sorted by date")
}

func TestCancelBooking(t *testing.T) {
	service := setupTestService()

	// Create a pending booking
	pendingBooking := &domain.Booking{
		ID:        "pending-booking",
		UserID:    "user1",
		ServiceID: "service1",
		Price:     40000,
		Status:    domain.StatusPending,
	}
	service.cache.SaveBooking(pendingBooking)
	service.mockRepository.SaveBooking(pendingBooking)

	// Create a confirmed booking
	confirmedBooking := &domain.Booking{
		ID:        "confirmed-booking",
		UserID:    "user2",
		ServiceID: "service2",
		Price:     60000,
		Status:    domain.StatusConfirmed,
	}
	service.cache.SaveBooking(confirmedBooking)
	service.mockRepository.SaveBooking(confirmedBooking)

	// Test canceling a pending booking
	err := service.CancelBooking(pendingBooking.ID)
	assert.NoError(t, err, "Expected no error when canceling a pending booking")

	// Test canceling a confirmed booking
	err = service.CancelBooking(confirmedBooking.ID)
	assert.Error(t, err, "Expected an error when canceling a confirmed booking")

	// Test canceling a non-existent booking
	err = service.CancelBooking("non-existent-booking")
	assert.Error(t, err, "Expected an error when canceling a non-existent booking")
}

func TestCancelExpiredBookings(t *testing.T) {
	service := setupTestService()

	// Create bookings with different creation times
	nonExpiredBooking := &domain.Booking{
		ID:        "non-expired-booking",
		UserID:    "user1",
		ServiceID: "service1",
		Price:     40000,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}
	expiredBooking := &domain.Booking{
		ID:        "expired-booking",
		UserID:    "user2",
		ServiceID: "service2",
		Price:     60000,
		Status:    domain.StatusPending,
		CreatedAt: time.Now().Add(-10 * time.Minute), // Set to expired
	}

	// Save bookings to cache
	service.cache.SaveBooking(nonExpiredBooking)
	service.cache.SaveBooking(expiredBooking)

	// Call CancelExpiredBookings
	service.CancelExpiredBookings()

	// Verify non-expired booking is still pending
	updatedNonExpiredBooking, _ := service.cache.GetBooking(nonExpiredBooking.ID)
	assert.Equal(t, domain.StatusPending, updatedNonExpiredBooking.Status, "Expected non-expired booking to remain pending")

	// Verify expired booking is canceled
	updatedExpiredBooking, _ := service.cache.GetBooking(expiredBooking.ID)
	assert.Equal(t, domain.StatusCanceled, updatedExpiredBooking.Status, "Expected expired booking to be canceled")
}
