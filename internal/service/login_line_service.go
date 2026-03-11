package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type LoginLineService interface {
	LoginLine(c fiber.Ctx, url string) error
	Callback(c fiber.Ctx, data string) error
}

type loginService struct{}

func NewLoginLineService() LoginLineService {
	return &loginService{}
}

const (
	channelID     = "2009413542"
	channelSecret = "58e22fbbdf081a462e044acb62e7397d"
	redirectURI   = "http://localhost:8080/callback"
	state         = "12345"
)

func (l *loginService) LoginLine(c fiber.Ctx, url string) error {
	return c.Redirect().Status(fiber.StatusFound).To(url)
}

func (l *loginService) Callback(c fiber.Ctx, data string) error {
	tokenURL := "https://api.line.me/oauth2/v2.1/token"

	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := io.ReadAll(resp.Body)

	var token map[string]interface{}
	json.Unmarshal(body, &token)

	accessToken := token["access_token"].(string)

	return getProfile(c, accessToken)
}

func getProfile(c fiber.Ctx, accessToken string) error {

	req, _ := http.NewRequest(
		"GET",
		"https://api.line.me/v2/profile",
		nil,
	)

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return c.SendString(err.Error())
	}

	body, _ := io.ReadAll(resp.Body)

	return c.Send(body)
}
