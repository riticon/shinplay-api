package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/user"
	"go.uber.org/zap"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthServiceIntr interface {
	// LoginOrSignupWithNumber(phoneNumber string, channel string) (token Token, err error)
	// LoginOrSignupWithEmail(email string, channel string) (token Token, err error)
	SendWhatsAppOTP(phoneNumber string) error
	GenerateOTP(phoneNumber string) (otp string, err error)
	VerifyWhatsAppOTP(phoneNumber, otp string) (bool, error)
	GenerateAuthTokens(user *ent.User) (token Token, err error)
	generateAccessToken(user *ent.User) (string, error)
	generateRefreshToken(user *ent.User) (string, error)
	// ValidateToken(token string) bool
	// RefreshToken(refreshToken string) (token Token, err error)
	// Logout(userId, sessionId string) error
}

type AuthService struct {
	userService *user.UserService
	otpService  *OTPService
	config      *config.Config
}

func NewAuthService(userService *user.UserService, otpService *OTPService, config *config.Config) *AuthService {
	return &AuthService{
		userService: userService,
		otpService:  otpService,
		config:      config,
	}
}

func (s *AuthService) GenerateAuthTokens(user *ent.User) (token Token, err error) {
	s.config.Logger.Info("Generating auth token for user", zap.Int("userId", user.ID))

	accessToken, _ := s.generateAccessToken(user)
	refreshToken, err := s.generateRefreshToken(user)

	token = Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return token, err
}

func (s *AuthService) generateAccessToken(user *ent.User) (string, error) {
	s.config.Logger.Info("Generating access token for user", zap.Int("userId", user.ID))
	claims := jwt.MapClaims{
		"sub": user.AuthID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *ent.User) (string, error) {
	s.config.Logger.Info("Generating refresh token for user", zap.Int("userId", user.ID))

	// valid for 30 days
	claims := jwt.MapClaims{
		"sub": user.AuthID,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) SendWhatsAppOTP(phoneNumber string) error {
	s.config.Logger.Info("Sending WhatsApp OTP", zap.String("phoneNumber", phoneNumber))
	url := "https://graph.facebook.com/v22.0/" + s.config.WhatsApp.PhoneId + "/messages"
	token := s.config.WhatsApp.Token
	otp, err := s.GenerateOTP(phoneNumber)

	s.config.Logger.Info("Generated OTP", zap.String("otp", otp))

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
						{"type": "text", "text": otp},
					},
				},
			},
		},
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		s.config.Logger.Error("Error marshalling payload", zap.Error(err))
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.config.Logger.Error("Error creating request", zap.Error(err))
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.config.Logger.Error("Error sending request", zap.Error(err))
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Print and handle response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", buf.String())

	if resp.StatusCode >= 400 {
		s.config.Logger.Error("Received error status from WhatsApp API", zap.String("status", resp.Status))
		return fmt.Errorf("received error status: %s", resp.Status)
	}

	return nil
}

func (s *AuthService) GenerateOTP(phoneNumber string) (string, error) {
	s.config.Logger.Info("Generating OTP for phone number", zap.String("phoneNumber", phoneNumber))

	// Find if user with phoneNumber exists
	// If not, create a new user with the phoneNumber
	// and return the OTP
	user, err := s.userService.FindOrCreateByPhone(phoneNumber)
	if err != nil {
		s.config.Logger.Error("Failed to find or create user", zap.Error(err))
		return "", fmt.Errorf("error finding or creating user: %w", err)
	}

	otp, err := s.otpService.CreateNewOTP(user)
	if err != nil {
		s.config.Logger.Error("Failed to create new OTP", zap.Error(err))
		return "", fmt.Errorf("error creating OTP: %w", err)
	}

	s.config.Logger.Info("OTP created successfully", zap.String("otp", otp.Otp))

	return otp.Otp, nil
}

type UserInfo struct {
	AuthID      string `json:"auth_id"`
	PhoneNumber string `json:"phone_number"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

func (s *AuthService) VerifyWhatsAppOTP(phoneNumber string, otp string) (Token, UserInfo, error) {
	s.config.Logger.Info("Verifying WhatsApp OTP", zap.String("phoneNumber", phoneNumber), zap.String("otp", otp))

	// Find user by phone number
	user, err := s.userService.FindByPhone(phoneNumber)
	if err != nil {
		s.config.Logger.Info("Failed to find user", zap.Error(err))
		return Token{}, UserInfo{}, fmt.Errorf("error finding user: %w", err)
	}

	// Check if OTP is valid
	is_valid, err := s.otpService.IsOTPValid(otp, user)
	if err != nil {
		s.config.Logger.Info("Failed to find OTP", zap.Error(err))
		return Token{}, UserInfo{}, fmt.Errorf("error finding OTP: %w", err)
	}

	if !is_valid {
		return Token{}, UserInfo{}, nil
	}

	// Expire the OTP
	s.otpService.ExpireOtp(otp, user)

	tokens, err := s.GenerateAuthTokens(user)

	s.config.Logger.Info("OTP is valid, generating tokens: ", zap.Any("tokens", "[Filtered]"))
	return tokens, UserInfo{
		AuthID:      user.AuthID,
		PhoneNumber: user.PhoneNumber,
		UserName:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}, err
}
