package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/touchsung/spd-fiber-booking-system/handler"
	"github.com/touchsung/spd-fiber-booking-system/middleware"
)

func SetupRoutes(app *fiber.App, bookingHandler *handler.BookingHandler) {
	// Add global middleware
	app.Use(middleware.RequestLogger())

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Booking routes
	app.Get("/bookings", bookingHandler.ListBookings)
	app.Post("/bookings", bookingHandler.Create)
	app.Get("/bookings/:id", bookingHandler.GetBooking)
	app.Delete("/bookings/:id", bookingHandler.CancelBooking)
}
