# Booking Service API

## Overview

This project is a Booking Service API built using the Go programming language and the Fiber web framework. It provides endpoints for managing bookings, including creating, listing, retrieving, and canceling bookings. The API also supports optional sorting and filtering of bookings.

## Features

- **Create Booking**: Allows users to create a new booking. A credit check is performed for bookings with a price greater than 50,000.
- **List Bookings**: Retrieve a list of all bookings with optional sorting by price or date and filtering for high-value bookings.
- **Get Booking by ID**: Fetch the details of a specific booking using its ID.
- **Cancel Booking**: Cancel a booking by its ID, with restrictions on canceling confirmed bookings.

## API Documentation

The API is documented using Swagger. You can access the Swagger UI at `/swagger/` once the application is running.

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/touchsung/spd-fiber-booking-system.git
   cd spd-fiber-booking-system
   ```

2. **Install dependencies**:
   Ensure you have Go installed, then run:

   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   go run cmd/api/main.go
   ```

## Usage

- **Base URL**: `http://localhost:3000`
- **API Docs**: `http://localhost:3000/swagger/`

### Endpoints

- **GET /bookings**: List all bookings with optional query parameters `sort` (price or date) and `high-value` (boolean).
- **POST /bookings**: Create a new booking. Requires a JSON body with `user_id`, `service_id`, and `price`.
- **GET /bookings/{id}**: Retrieve a booking by its ID.
- **DELETE /bookings/{id}**: Cancel a booking by its ID.

## Development

### Running Tests

To run the tests, use the following command:

```bash
go test ./...
```

### Code Structure

- **cmd**: Contains the main entry point for the application.
- **dto**: Data Transfer Objects used in the application.
- **handler**: Contains the HTTP handlers for the API endpoints.
- **middleware**: Middleware components for the application.
- **models**: Defines the domain models.
- **repository**: Repository layer for data access.
- **router**: Defines the routes for the application.
- **usecase**: Contains the business logic for managing bookings.
- **utils**: Utility functions.
- **docs**: Documentation files.
