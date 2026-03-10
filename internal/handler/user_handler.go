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
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := validator.Validate(req); len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.service.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, "Email already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response.JSON(w, http.StatusCreated, user)
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
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	response.JSON(w, http.StatusOK, user)
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
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := validator.Validate(req); len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, "Email already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.JSON(w, http.StatusOK, user)
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
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
