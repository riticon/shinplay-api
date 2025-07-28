package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(prefix string, router fiber.Router) {
	auth := router.Group(prefix)

	auth.Get("/sms/send-otp", sendSMSOTP)
	auth.Post("/whatsapp/send-otp", sendWhatsAppOTP)
}

// POST /auth/sms-otp
func sendSMSOTP(ctx *fiber.Ctx) error {
	// authService AuthAuthService = auth.NewAuthService()
	// authService.S
	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":  "success",
			"message": "OTP sent successfully",
		})
}

// POST /auth/whatsapp/send-otp
// This endpoint sends an OTP via WhatsApp to the provided phone number.
// It expects a JSON body with the { phone_number }

type WhatsAppOtpBody struct {
	PhoneNumber string `json:"phone_number" xml:"phone_number" form:"phone_number"`
}

func sendWhatsAppOTP(ctx *fiber.Ctx) error {
	authService := NewAuthService()

	body := new(WhatsAppOtpBody)
	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Please provide a valid phone number",
		})
	}

	go authService.SendWhatsAppOTP(body.PhoneNumber, "123456") // Example OTP

	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":  "success",
			"message": "WhatsApp OTP sent successfully",
		})
}
