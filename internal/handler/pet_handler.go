package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/maythitirat/pet-log-api/internal/model"
	"github.com/maythitirat/pet-log-api/internal/service"
	"github.com/maythitirat/pet-log-api/pkg/validator"
)

// PetHandler handles pet-related HTTP requests
type PetHandler struct {
	service service.PetService
}

// NewPetHandler creates a new pet handler
func NewPetHandler(svc service.PetService) *PetHandler {
	return &PetHandler{service: svc}
}

// Create handles POST /pets
// @Summary Create a new pet
// @Description Create a new pet with the provided information
// @Tags Pets
// @Accept json
// @Produce json
// @Param pet body model.CreatePetRequest true "Pet data"
// @Success 201 {object} model.PetResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pets [post]
func (h *PetHandler) Create(c fiber.Ctx) error {
	var req model.CreatePetRequest
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

	pet, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create pet"})
	}

	return c.Status(fiber.StatusCreated).JSON(pet)
}

// GetByID handles GET /pets/{id}
// @Summary Get a pet by ID
// @Description Get a pet's information by their ID
// @Tags Pets
// @Produce json
// @Param id path int true "Pet ID"
// @Success 200 {object} model.PetResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pets/{id} [get]
func (h *PetHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pet ID"})
	}

	pet, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get pet"})
	}

	return c.JSON(pet)
}

// GetAll handles GET /pets
// @Summary Get all pets
// @Description Get a list of all pets with pagination
// @Tags Pets
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {array} model.PetResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pets [get]
func (h *PetHandler) GetAll(c fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	pets, err := h.service.GetAll(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get pets"})
	}

	return c.JSON(pets)
}

// GetByOwnerID handles GET /users/{userId}/pets
// @Summary Get pets by owner ID
// @Description Get all pets belonging to a specific owner
// @Tags Pets
// @Produce json
// @Param userId path int true "Owner (User) ID"
// @Success 200 {array} model.PetResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/pets [get]
func (h *PetHandler) GetByOwnerID(c fiber.Ctx) error {
	ownerID, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid owner ID"})
	}

	pets, err := h.service.GetByOwnerID(c.Context(), ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get pets"})
	}

	return c.JSON(pets)
}

// Update handles PUT /pets/{id}
// @Summary Update a pet
// @Description Update a pet's information
// @Tags Pets
// @Accept json
// @Produce json
// @Param id path int true "Pet ID"
// @Param pet body model.UpdatePetRequest true "Pet data to update"
// @Success 200 {object} model.PetResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pets/{id} [put]
func (h *PetHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pet ID"})
	}

	var req model.UpdatePetRequest
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

	pet, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update pet"})
	}

	return c.JSON(pet)
}

// Delete handles DELETE /pets/{id}
// @Summary Delete a pet
// @Description Delete a pet by their ID
// @Tags Pets
// @Produce json
// @Param id path int true "Pet ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pets/{id} [delete]
func (h *PetHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pet ID"})
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete pet"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
