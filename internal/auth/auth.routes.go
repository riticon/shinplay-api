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

type WhatsAppNumberBody struct {
	PhoneNumber string `json:"phone_number" xml:"phone_number" form:"phone_number"`
}

func sendWhatsAppOTP(ctx *fiber.Ctx) error {
	authService := NewAuthService()
	body := new(WhatsAppNumberBody)
	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request",
		})
	}

	phoneNumber := body.PhoneNumber

	println("Received phone number:", phoneNumber)

	authService.SendWhatsAppOTP(phoneNumber, "123456") // Example OTP
	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":  "success",
			"message": "WhatsApp OTP sent successfully",
		})
}
