package v1

import (
	"context"
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/validator"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/gofiber/fiber/v2"
)

type AdminService interface {
	Create(ctx context.Context, createAdmin *dto.CreateAdminValue) (*entities.Admin, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Admin, error)
	Update(ctx context.Context, admin *dto.UpdateAdminValue) (*entities.Admin, error)
	Delete(ctx context.Context, uuid string) error
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Admin, error)
}

type AdminHandler struct {
	adminService AdminService
	validator    *validator.Validator
}

func NewAdminHandler(bizkitEduApp *app.BizkitEduApp) *AdminHandler {
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	adminService := services.NewAdminService(adminStorage)

	return &AdminHandler{
		adminService: adminService,
		validator:    bizkitEduApp.Validator,
	}
}

func (h AdminHandler) create(c *fiber.Ctx) error {
	var createAdmin dto.CreateAdminValue

	if err := c.BodyParser(&createAdmin); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(createAdmin)
	if errValidate != nil {
		return errValidate
	}

	adminValueObject, err := h.adminService.Create(c.Context(), &createAdmin)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   adminValueObject,
	})

}

func (h AdminHandler) getAll(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	values, err := h.adminService.GetAll(c.Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"status": true,
		"body":   values,
	})
}

func (h AdminHandler) update(c *fiber.Ctx) error {
	var updateAdmin dto.UpdateAdminValue

	if err := c.BodyParser(&updateAdmin); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(&updateAdmin)
	if errValidate != nil {
		return errValidate
	}

	updatedValue, err := h.adminService.Update(c.Context(), &updateAdmin)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   updatedValue,
	})
}

func (h AdminHandler) Setup(router fiber.Router, handler fiber.Handler) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(handler)
	adminGroup.Post("/create", h.create, handler)
	adminGroup.Get("/all", h.getAll, handler)
	adminGroup.Put("/update", h.update, handler)
}
