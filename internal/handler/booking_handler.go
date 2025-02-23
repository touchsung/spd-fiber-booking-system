package handler

import (
	"github.com/touchsung/spd-fiber-booking-system/internal/domain"
	"github.com/touchsung/spd-fiber-booking-system/internal/service"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// Create godoc
// @Summary Create a new booking
// @Description Create a new booking with the provided details. A credit check is performed for bookings with a price greater than 50,000.
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body domain.BookingRequest true "Booking Request"
// @Success 201 {object} domain.Booking
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) Create(c *fiber.Ctx) error {
	var request domain.BookingRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	booking, err := h.bookingService.CreateBooking(request)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create booking",
		})
	}

	return c.Status(201).JSON(booking)
}

// GetBooking godoc
// @Summary Get a booking by ID
// @Description Get a booking's details by its ID. The booking is retrieved from cache first, then from the mock repository if not found.
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} domain.Booking
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")

	booking, err := h.bookingService.GetBooking(bookingID)
	if err != nil {
		return c.Status(404).JSON(ErrorResponse{
			Error: "Booking not found",
		})
	}

	return c.JSON(booking)
}

// ListBookings godoc
// @Summary List all bookings
// @Description Get a list of all bookings with optional sorting and filtering. Sort by price or date, or default to ID. Filter high-value bookings (price > 50,000).
// @Tags bookings
// @Accept json
// @Produce json
// @Param sort query string false "Sort by field (price or date)"
// @Param high-value query bool false "Filter high-value bookings (price > 50,000)"
// @Success 200 {array} domain.Booking
// @Router /bookings [get]
func (h *BookingHandler) ListBookings(c *fiber.Ctx) error {
	// Parse query parameters
	var sortBy *domain.SortOption
	if sort := c.Query("sort"); sort != "" {
		switch sort {
		case "price":
			option := domain.SortByPrice
			sortBy = &option
		case "date":
			option := domain.SortByDate
			sortBy = &option
		}
	}

	var highValueOnly *bool
	if highValue := c.Query("high-value"); highValue != "" {
		value := highValue == "true"
		highValueOnly = &value
	}

	bookings := h.bookingService.ListBookings(sortBy, highValueOnly)
	return c.JSON(bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking by its ID. Cannot cancel confirmed bookings. The booking is checked in both cache and repository.
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse "Cannot cancel confirmed booking"
// @Failure 404 {object} ErrorResponse "Booking not found"
// @Router /bookings/{id} [delete]
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")

	err := h.bookingService.CancelBooking(bookingID)
	if err != nil {
		if err.Error() == "cannot cancel a confirmed booking" {
			return c.Status(400).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(404).JSON(ErrorResponse{
			Error: "Booking not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Booking canceled successfully",
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}
