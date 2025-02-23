package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/touchsung/spd-fiber-booking-system/docs" // This will be generated
	"github.com/touchsung/spd-fiber-booking-system/internal/handler"
	"github.com/touchsung/spd-fiber-booking-system/internal/middleware"
	"github.com/touchsung/spd-fiber-booking-system/internal/repository"
	"github.com/touchsung/spd-fiber-booking-system/internal/service"
)

// @title Booking API
// @version 1.0
// @description This is a booking service API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email your-email@domain.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	// Add global middleware
	app.Use(middleware.RequestLogger())

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize dependencies
	cache := repository.NewInMemoryCache()
	mockRepo := repository.NewMockRepository()
	bookingService := service.NewBookingService(cache, mockRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Routes
	app.Get("/bookings", bookingHandler.ListBookings)
	app.Post("/bookings", bookingHandler.Create)
	app.Get("/bookings/:id", bookingHandler.GetBooking)
	app.Delete("/bookings/:id", bookingHandler.CancelBooking)

	// Start background task
	runBackgroundTask(bookingService)

	app.Listen(":3000")
}

// Function to run background task for checking expired bookings
func runBackgroundTask(bookingService *service.BookingService) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				bookingService.CancelExpiredBookings()
			}
		}
	}()
}
