package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/auth/otp"
	"github.com/shinplay/internal/auth/session"
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
	GoogleOauthSignIn(idToken string, ipAddress string, userAgent string) (token Token, userInfo UserInfo, sessionID string, err error)
	VerifyWhatsAppOTP(phoneNumber, otp string) (bool, error)
	GenerateAuthTokens(user *ent.User) (token Token, err error)
	generateAccessToken(user *ent.User) (string, error)
	generateRefreshToken(user *ent.User) (string, error)
	LoginUser(user *ent.User) (token Token, err error)
	ValidateToken(token string) bool
	// RefreshToken(refreshToken string) (token Token, err error)
	// Logout(userId, sessionId string) error
}

type AuthService struct {
	userService       *user.UserService
	otpService        *otp.OTPService
	sessionRepository *session.SessionRepository
	config            *config.Config
}

type UserInfo struct {
	AuthID      string `json:"auth_id"`
	PhoneNumber string `json:"phone_number"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

func NewAuthService(userService *user.UserService, otpService *otp.OTPService, sessionRepository *session.SessionRepository, config *config.Config) *AuthService {
	return &AuthService{
		userService:       userService,
		otpService:        otpService,
		sessionRepository: sessionRepository,
		config:            config,
	}
}

func (s *AuthService) LoginUser(user *ent.User, ipAddress string, userAgent string) (Token, UserInfo, string, error) {
	// Create a new session for the user
	tokens, err := s.GenerateAuthTokens(user)
	if err != nil {
		s.config.Logger.Error("Failed to generate auth tokens", zap.Error(err))
		return Token{}, UserInfo{}, "", err
	}

	userInfo := UserInfo{
		AuthID:      user.AuthID,
		PhoneNumber: user.PhoneNumber,
		UserName:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}

	s.config.Logger.Info("Creating Session", zap.Any("id", ipAddress), zap.Any("userAgent", userAgent))

	session, err := s.sessionRepository.CreateNewSession(
		context.Background(),
		user,
		tokens.RefreshToken,
		time.Now().Add(30*24*time.Hour), // 30 days expiration
		userAgent,
		ipAddress,
	)

	if err != nil {
		s.config.Logger.Error("Failed to create session", zap.Error(err))
		return Token{}, UserInfo{}, "", err
	}

	return tokens, userInfo, session.SessionID, nil
}

func (s *AuthService) GoogleOauthSignIn(idToken string, ipAddress string, userAgent string) (user *ent.User, err error) {
	// Validate the ID token
	payload, err := idtoken.Validate(context.Background(), idToken, s.config.Google.ClientID)

	if err != nil {
		s.config.Logger.Error("Failed to validate ID token", zap.Error(err))
		return nil, fmt.Errorf("failed to validate ID token: %w", err)
	}

	// Find user by email
	email := payload.Claims["email"].(string)
	user, err = s.userService.FindOrCreateByEmail(email)

	if err != nil {
		s.config.Logger.Error("Failed to find or create user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) GenerateAuthTokens(user *ent.User) (token Token, err error) {
	accessToken, _ := s.generateAccessToken(user)
	refreshToken, err := s.generateRefreshToken(user)

	token = Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return token, err
}

func (s *AuthService) generateAccessToken(user *ent.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.AuthID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *ent.User) (string, error) {
	// valid for 30 days
	claims := jwt.MapClaims{
		"sub": user.AuthID,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) SendWhatsAppOTP(phoneNumber string) error {
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

	return otp.Otp, nil
}

func (s *AuthService) VerifyWhatsAppOTP(phoneNumber string, otp string) (bool, *ent.User) {
	// Find user by phone number
	user, err := s.userService.FindByPhone(phoneNumber)
	if err != nil {
		s.config.Logger.Info("Failed to find user", zap.Error(err))
		return false, nil
	}

	// Check if OTP is valid
	is_valid, err := s.otpService.IsOTPValid(otp, user)
	if err != nil {
		s.config.Logger.Info("Failed to validate OTP", zap.Error(err))
		return false, nil
	}

	// Expire the OTP
	s.otpService.ExpireOtp(otp, user.ID)

	return is_valid, user
}

func (s *AuthService) ValidateToken(token string) (bool, *ent.User) {
	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		s.config.Logger.Error("Failed to parse token", zap.Error(err))
		return false, nil
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		s.config.Logger.Error("Invalid token claims")
		return false, nil
	}

	sub, err := claims.GetSubject()
	if err != nil {
		s.config.Logger.Error("Failed to get subject from claims", zap.Error(err))
		return false, nil
	}

	s.config.Logger.Info("Token validated successfully", zap.Any("claims", sub))

	user, err := s.userService.FindUserByAuthID(sub)
	if err != nil {
		s.config.Logger.Error("Failed to find user by auth ID", zap.String("auth_id", sub), zap.Error(err))
		return false, nil
	}

	return parsedToken.Valid, user
}
