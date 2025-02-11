package auth

import (
	"github.com/gofiber/fiber/v2"
)

// AuthHandler menangani HTTP request untuk auth
type AuthHandler struct {
	authService AuthService
}

// NewAuthHandler membuat instance handler baru
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

// Register menangani registrasi user
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request RegisterRequest

	// Parse request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	user, err := request.ConvertToUser()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Register user
	createdUser, err := h.authService.RegisterUser(*user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Response sukses tanpa password
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"username":     createdUser.Username,
			"email":        createdUser.Email,
			"phone_number": createdUser.PhoneNumber,
			"address":      createdUser.Address,
			"dob":          createdUser.DOB.Format("2006-01-02"),
		},
	})
}
