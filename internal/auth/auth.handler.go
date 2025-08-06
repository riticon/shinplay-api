package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/auth/session"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *AuthService
	config      *config.Config
}

func NewAuthHandler(authService *AuthService, sessionRepository *session.SessionRepository, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

type WhatsAppOtpBody struct {
	PhoneNumber string `json:"phone_number" xml:"phone_number" form:"phone_number"`
}

func (h *AuthHandler) SendWhatsAppOTP(ctx *fiber.Ctx) error {
	body := new(WhatsAppOtpBody)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Please provide a valid phone number",
		})
	}

	go h.authService.SendWhatsAppOTP(body.PhoneNumber) // Example OTP

	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":  "success",
			"message": "WhatsApp OTP sent successfully",
		})
}

type OTPBody struct {
	Otp         string `json:"otp" xml:"otp" form:"otp"`
	PhoneNumber string `json:"phone_number" xml:"phone_number" form:"phone_number"`
}

func (h *AuthHandler) VerifyWhatsAppOTP(ctx *fiber.Ctx) error {
	body := new(OTPBody)

	if err := ctx.BodyParser(body); err != nil {
		h.config.Logger.Warn("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Please provide a valid OTP",
		})
	}

	isValid, user := h.authService.VerifyWhatsAppOTP(body.PhoneNumber, body.Otp)

	if !isValid {
		h.config.Logger.Info("Failed to verify WhatsApp OTP", zap.String("phone_number", body.PhoneNumber))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid OTP or OTP expired",
		})
	}

	h.config.Logger.Info("User verified successfully", zap.Any("user_id", user))

	tokens, userInfo, sessionId, err := h.authService.LoginUser(user, ctx.IP(), ctx.Get("User-Agent"))

	if err != nil {
		h.config.Logger.Error("Failed to create session", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to Login, please try again later",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "OTP verified successfully",
		"data": fiber.Map{
			"access_token": tokens.AccessToken,
			"user":         userInfo,
		},
	})
}

func (h *AuthHandler) GoogleOauthSignin(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing Id Token in Authorization header",
		})
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := h.authService.GoogleOauthSignIn(idToken, ctx.IP(), ctx.Get("User-Agent"))

	if err != nil {
		h.config.Logger.Error("Failed to sign in with Google", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to sign in with Google, please try again later",
		})
	}

	tokens, userInfo, sessionId, err := h.authService.LoginUser(user, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		h.config.Logger.Error("Failed to create session", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to Login, please try again later",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "User signed in successfully",
		"data": fiber.Map{
			"access_token": tokens.AccessToken,
			"user":         userInfo,
		},
	})
}

func (h *AuthHandler) AuthenticateUser(ctx *fiber.Ctx) error {
	h.config.Logger.Info("Authenticating user")

	header := ctx.Get("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing Authorization header",
		})
	}

	isValid, user := h.authService.ValidateToken(token)
	if !isValid {
		h.config.Logger.Info("Invalid or expired token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or expired token",
		})
	}

	ctx.Locals("user", user)

	return ctx.Next()
}
