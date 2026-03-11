package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/maythitirat/pet-log-api/internal/config"
	"github.com/maythitirat/pet-log-api/internal/service"
)

type LoginLineHandler struct {
	service service.LoginLineService
	cfg     config.LineConfig
}

func NewLoginLineHandler(svc service.LoginLineService, config *config.LineConfig) *LoginLineHandler {
	return &LoginLineHandler{service: svc, cfg: *config}
}

func (h *LoginLineHandler) LoginLine(c fiber.Ctx) error {
	url := fmt.Sprintf(
		"https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=profile%%20openid&state=%s",
		h.cfg.ChannelID,
		h.cfg.RedirectURI,
		h.cfg.State,
	)

	return h.service.LoginLine(c, url)
}

func (h *LoginLineHandler) Callback(c fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if state != "12345" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid state"})
	}

	data := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code,
		h.cfg.RedirectURI,
		h.cfg.ChannelID,
		h.cfg.ChannelSecret,
	)

	return h.service.Callback(c, data)
}
