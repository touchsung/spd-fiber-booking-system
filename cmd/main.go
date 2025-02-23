package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/touchsung/spd-fiber-booking-system/docs" // This will be generated
	"github.com/touchsung/spd-fiber-booking-system/handler"
	"github.com/touchsung/spd-fiber-booking-system/repository"
	"github.com/touchsung/spd-fiber-booking-system/router"
	"github.com/touchsung/spd-fiber-booking-system/usecase"
	"github.com/touchsung/spd-fiber-booking-system/utils"
)

// @title Booking API
// @version 1.0
// @description This is a booking service API
// @termsOfService http://swagger.io/terms/

// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	// Initialize dependencies
	cache := utils.NewInMemoryCache()
	mockRepo := repository.NewMockRepository()
	bookingService := usecase.NewBookingService(cache, mockRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Setup routes
	router.SetupRoutes(app, bookingHandler)

	// Start background task
	runBackgroundTask(bookingService)

	app.Listen(":3000")
}

// Function to run background task for checking expired bookings
func runBackgroundTask(bookingService *usecase.BookingService) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			bookingService.CancelExpiredBookings()
		}
	}()
}
