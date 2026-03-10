package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/maythitirat/pet-log-api/internal/model"
	"github.com/maythitirat/pet-log-api/internal/service"
	"github.com/maythitirat/pet-log-api/pkg/response"
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
func (h *PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := validator.Validate(req); len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	pet, err := h.service.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create pet")
		return
	}

	response.JSON(w, http.StatusCreated, pet)
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
func (h *PetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	pet, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			response.Error(w, http.StatusNotFound, "Pet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to get pet")
		return
	}

	response.JSON(w, http.StatusOK, pet)
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
func (h *PetHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	pets, err := h.service.GetAll(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get pets")
		return
	}

	response.JSON(w, http.StatusOK, pets)
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
func (h *PetHandler) GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}

	pets, err := h.service.GetByOwnerID(r.Context(), ownerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get pets")
		return
	}

	response.JSON(w, http.StatusOK, pets)
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
func (h *PetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	var req model.UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := validator.Validate(req); len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	pet, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			response.Error(w, http.StatusNotFound, "Pet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update pet")
		return
	}

	response.JSON(w, http.StatusOK, pet)
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
func (h *PetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPetNotFound) {
			response.Error(w, http.StatusNotFound, "Pet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete pet")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
