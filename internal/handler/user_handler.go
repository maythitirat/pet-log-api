package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/maythitirat/pet-log-api/internal/model"
	"github.com/maythitirat/pet-log-api/internal/service"
	"github.com/maythitirat/pet-log-api/pkg/validator"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// Create handles POST /users
// @Summary Create a new user
// @Description Create a new user (registration)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.CreateUserRequest true "User data"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate request
	if validationErrors := validator.Validate(req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation Error",
			"message": "One or more fields failed validation",
			"details": validationErrors,
		})
	}

	user, err := h.service.Create(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetByID handles GET /users/{id}
// @Summary Get a user by ID
// @Description Get a user's information by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user"})
	}

	return c.JSON(user)
}

// Update handles PUT /users/{id}
// @Summary Update a user
// @Description Update a user's information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.UpdateUserRequest true "User data to update"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req model.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate request
	if validationErrors := validator.Validate(req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation Error",
			"message": "One or more fields failed validation",
			"details": validationErrors,
		})
	}

	user, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(user)
}

// Delete handles DELETE /users/{id}
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
