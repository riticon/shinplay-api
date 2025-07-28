package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/user"
	"go.uber.org/zap"
)

type Token struct {
	accessToken  string
	refreshToken string
}

type AuthService interface {
	LoginOrSignupWithNumber(phoneNumber string, channel string) (token Token, err error)
	LoginOrSignupWithEmail(email string, channel string) (token Token, err error)
	SendWhatsAppOTP(phoneNumber string) error
	GenerateOTP(phoneNumber string) (otp string, err error)
	VerifyWhatsAppOTP(phoneNumber, otp string) (token Token, err error)
	ValidateToken(token string) bool
	RefreshToken(refreshToken string) (token Token, err error)
	Logout(userId, sessionId string) error
}

type AuthServiceImpl struct {
	userService user.UserService
	config      *config.Config
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{
		userService: user.NewUserService(),
		config:      config.GetConfig(),
	}
}

func SendSMSOTP(n *AuthService) {
	config.GetConfig().Logger.Info("Sending OTP")
}

// func (s *AuthServiceImpl) SendWhatsAppOtp(phoneNumber string, channel string) error {}

func (s *AuthServiceImpl) SendWhatsAppOTP(phoneNumber, otp string) error {
	url := "https://graph.facebook.com/v22.0/" + s.config.WhatsApp.PhoneId + "/messages"
	token := s.config.WhatsApp.Token
	otp, err := s.GenerateOTP(phoneNumber)

	if err != nil {
		return fmt.Errorf("error generating OTP: %w", err)
	}

	// Payload
	payload := map[string]any{
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

func (s *AuthServiceImpl) GenerateOTP(phoneNumber string) (string, error) {
	// Generate a random OTP (for simplicity, using a static value here)

	// Find if user with phoneNumber exists
	// If not, create a new user with the phoneNumber
	// and return the OTPW
	user, err := s.userService.FindOrCreateByPhone(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("error finding or creating user: %w", err)
	}

	s.config.Logger.Info("User found or created", zap.Int("user_id", user.ID))

	otp := "123456"
	return otp, nil
}
