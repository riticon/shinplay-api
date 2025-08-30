package country

import "github.com/gofiber/fiber/v2"

type Handler struct {
	Svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{Svc: svc}
}

// GET /countries  (public)
func (h *Handler) List(ctx *fiber.Ctx) error {
	list, err := h.Svc.Load()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load countries",
			"detail": err.Error(), // optional: remove in prod if you prefer
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(list)
}
