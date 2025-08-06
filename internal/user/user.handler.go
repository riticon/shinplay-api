package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	CheckUsernameAvailability(ctx *fiber.Ctx) (bool, error)
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
