package models

import (
	"time"
)

type CreditCheckResult struct {
	BookingID string
	Status    BookingStatus
}

type SortOption string

// @Description Sort option enum
const (
	SortByPrice SortOption = "price"
	SortByDate  SortOption = "date"
)

// BookingStatus represents the status of a booking
// @Description Booking status enum
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"   // Initial status
	StatusConfirmed BookingStatus = "confirmed" // Booking is confirmed
	StatusRejected  BookingStatus = "rejected"  // Booking is rejected
	StatusCanceled  BookingStatus = "canceled"  // Booking is canceled
)

// Booking represents a booking record
// @Description Booking information
type Booking struct {
	ID        string        `json:"id" example:"20240319123456"`
	UserID    string        `json:"user_id" example:"user123"`
	ServiceID string        `json:"service_id" example:"service456"`
	Price     float64       `json:"price" example:"60000"`
	Status    BookingStatus `json:"status" example:"pending"`
	CreatedAt time.Time     `json:"created_at"`
}

// BookingRequest represents the incoming booking request
// @Description Booking creation request
type BookingRequest struct {
	UserID    string  `json:"user_id" example:"user123" validate:"required"`
	ServiceID string  `json:"service_id" example:"service456" validate:"required"`
	Price     float64 `json:"price" example:"60000" validate:"required,gt=0"`
}
