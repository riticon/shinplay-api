package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	CheckUsernameAvailability(ctx *fiber.Ctx) (bool, error)
	ChangeUsername(ctx *fiber.Ctx) error
}

type UserHandler struct {
	config      *config.Config
	userService *UserService
}

func NewUserHandler(userService *UserService, config *config.Config) *UserHandler {
	return &UserHandler{
		config:      config,
		userService: userService,
	}
}

type UsernameQuery struct {
	Check string `json:"check" xml:"check" form:"check" query:"check" param:"check"`
}

func (h *UserHandler) CheckUsernameAvailability(ctx *fiber.Ctx) error {
	var query = new(UsernameQuery)
	if err := ctx.QueryParser(query); err != nil {
		return err
	}

	h.config.Logger.Info("Checking username availability", zap.String("username", query.Check))
	_, err := h.userService.FindByUsername(query.Check)
	if err != nil {
		if ent.IsNotFound(err) {
			return ctx.JSON(fiber.Map{"available": true})
		}
		return err
	}

	return ctx.JSON(fiber.Map{"available": false})
}

type ChangeUsernameBody struct {
	NewUsername string `json:"new_username" xml:"new_username" form:"new_username"`
}

func (h *UserHandler) ChangeUsername(ctx *fiber.Ctx) error {
	var body = new(ChangeUsernameBody)
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	currentUser := ctx.Locals("user").(*ent.User)

	h.config.Logger.Info("Changing username", zap.String("authId", currentUser.AuthID), zap.String("newUsername", body.NewUsername))
	_, usernameTaken, err := h.userService.ChangeUsername(currentUser.AuthID, body.NewUsername)
	if usernameTaken {
		h.config.Logger.Warn("Username already taken", zap.String("newUsername", body.NewUsername))
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already taken"})
	}

	if err != nil {
		h.config.Logger.Error("Failed to change username", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to change username"})
	}

	return ctx.JSON(fiber.Map{"message": "Username changed successfully"})
}
