package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	CheckUsernameAvailability(username string) (bool, error)
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

type UsernameParam struct {
	check string `json:"username" xml:"username" form:"username"`
}

func (h *UserHandler) CheckUsernameAvailability(ctx *fiber.Ctx) error {
	var params UsernameParam
	if err := ctx.ParamsParser(&params); err != nil {
		return err
	}

	h.config.Logger.Info("Checking username availability", zap.String("username", params.check))
	_, err := h.userService.FindByUsername(params.check)
	if err != nil {
		if ent.IsNotFound(err) {
			return ctx.JSON(fiber.Map{"available": true})
		}
		return err
	}

	return ctx.JSON(fiber.Map{"available": false})
}
