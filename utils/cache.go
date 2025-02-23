package utils

import (
	"sync"

	"github.com/touchsung/spd-fiber-booking-system/models"
)

type InMemoryCache struct {
	bookings map[string]*models.Booking
	mutex    sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		bookings: make(map[string]*models.Booking),
	}
}

func (c *InMemoryCache) SaveBooking(booking *models.Booking) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.bookings[booking.ID] = booking
}

func (c *InMemoryCache) UpdateBookingStatus(bookingID string, status models.BookingStatus) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if booking, exists := c.bookings[bookingID]; exists {
		booking.Status = status
	}
}

func (c *InMemoryCache) GetBooking(bookingID string) (*models.Booking, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	booking, exists := c.bookings[bookingID]
	return booking, exists
}

func (c *InMemoryCache) GetAllBookings() []*models.Booking {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	bookings := make([]*models.Booking, 0, len(c.bookings))
	for _, booking := range c.bookings {
		bookings = append(bookings, booking)
	}
	return bookings
}

func (c *InMemoryCache) DeleteBooking(bookingID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.bookings, bookingID)
}
