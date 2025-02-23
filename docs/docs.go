// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "your-email@domain.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/bookings": {
            "get": {
                "description": "Get a list of all bookings with optional sorting and filtering. Sort by price or date, or default to ID. Filter high-value bookings (price \u003e 50,000).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "List all bookings",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Sort by field (price or date)",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter high-value bookings (price \u003e 50,000)",
                        "name": "high-value",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Booking"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new booking with the provided details. A credit check is performed for bookings with a price greater than 50,000.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Create a new booking",
                "parameters": [
                    {
                        "description": "Booking Request",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.BookingRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Booking"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/bookings/{id}": {
            "get": {
                "description": "Get a booking's details by its ID. The booking is retrieved from cache first, then from the mock repository if not found.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get a booking by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Booking ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Booking"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Cancel a booking by its ID. Cannot cancel confirmed bookings. The booking is checked in both cache and repository.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Cancel a booking",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Booking ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Cannot cancel confirmed booking",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Booking not found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Booking": {
            "description": "Booking information",
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string",
                    "example": "20240319123456"
                },
                "price": {
                    "type": "number",
                    "example": 60000
                },
                "service_id": {
                    "type": "string",
                    "example": "service456"
                },
                "status": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.BookingStatus"
                        }
                    ],
                    "example": "pending"
                },
                "user_id": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "domain.BookingRequest": {
            "description": "Booking creation request",
            "type": "object",
            "required": [
                "price",
                "service_id",
                "user_id"
            ],
            "properties": {
                "price": {
                    "type": "number",
                    "example": 60000
                },
                "service_id": {
                    "type": "string",
                    "example": "service456"
                },
                "user_id": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "domain.BookingStatus": {
            "description": "Booking status enum",
            "type": "string",
            "enum": [
                "pending",
                "confirmed",
                "rejected",
                "canceled"
            ],
            "x-enum-comments": {
                "StatusCanceled": "Booking is canceled",
                "StatusConfirmed": "Booking is confirmed",
                "StatusPending": "Initial status",
                "StatusRejected": "Booking is rejected"
            },
            "x-enum-varnames": [
                "StatusPending",
                "StatusConfirmed",
                "StatusRejected",
                "StatusCanceled"
            ]
        },
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Booking API",
	Description:      "This is a booking service API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
