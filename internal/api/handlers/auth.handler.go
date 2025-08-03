package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *auth.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *auth.AuthService, config *config.Config) *AuthHandler {
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Please provide a valid OTP",
		})
	}

	token, userInfo, err := h.authService.VerifyWhatsAppOTP(body.PhoneNumber, body.Otp)

	if err != nil {
		h.config.Logger.Warn("Failed to verify WhatsApp OTP", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid OTP or OTP expired",
		})
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":  "success",
			"message": "OTP verified successfully",
			"data": fiber.Map{
				"access_token": token.AccessToken,
				"user":         userInfo,
			},
		})
}
