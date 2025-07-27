package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shinplay/internal/config"
)

type Token struct {
	accessToken  string
	refreshToken string
}

type AuthService interface {
	LoginOrSignupWithNumber(phoneNumber string, channel string) (token Token, err error)
	LoginOrSignupWithEmail(email string, channel string) (token Token, err error)
	SendWhatsAppOTP(phoneNumber string) error
	ValidateToken(token string) bool
	RefreshToken(refreshToken string) (token Token, err error)
	Logout(userId, sessionId string) error
}

type AuthServiceImpl struct {
	config *config.Config
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{
		config: config.GetConfig(),
	}
}

func SendSMSOTP(n *AuthService) {
	config.GetConfig().Logger.Info("Sending OTP")
}

// func (s *AuthServiceImpl) SendWhatsAppOtp(phoneNumber string, channel string) error {}

func (s *AuthServiceImpl) SendWhatsAppOTP(phoneNumber, otp string) error {
	url := "https://graph.facebook.com/v22.0/" + s.config.WhatsApp.PhoneId + "/messages"
	token := s.config.WhatsApp.Token

	println("Sending WhatsApp OTP to:", phoneNumber)

	// Payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                phoneNumber,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "otp_login",
			"language": map[string]string{
				"code": "en_US",
			},
			"components": []interface{}{
				map[string]interface{}{
					"type": "body",
					"parameters": []map[string]string{
						{"type": "text", "text": otp},
						{"type": "text", "text": "+91 7019331704"},
					},
				},
				map[string]interface{}{
					"type":     "button",
					"sub_type": "url",
					"index":    "0",
					"parameters": []map[string]string{
						{"type": "text", "text": "abc123"},
					},
				},
			},
		},
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Print and handle response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", buf.String())

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received error status: %s", resp.Status)
	}

	return nil

}
