package v1

import (
	"context"
	"github.com/Louffty/green-code-moscow/cmd/app"
	apiDto "github.com/Louffty/green-code-moscow/internal/adapters/controller/api/dto"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/validator"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/Louffty/green-code-moscow/internal/domain/usecases/event"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/gofiber/fiber/v2"
)

type EventsUserService interface {
	Create(ctx context.Context, eventsUser *dto.CreateEventsUser, userUUID string) (*entities.EventsUser, error)
	GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) ([]*entities.EventsUser, error)
	GetAllByUser(ctx context.Context, userUUID string, limit, offset int) ([]*entities.EventsUser, error)
	Delete(ctx context.Context, UUID string) error
}

type EventsUserUseCase interface {
	GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) (*dto.ReturnEventsUser, error)
	GetAllByUser(ctx context.Context, userUUID string, limit, offset int) (*dto.ReturnUsersEvents, error)
}

type EventsUserHandler struct {
	eventsUserService EventsUserService
	eventsUserUseCase EventsUserUseCase
	validator         *validator.Validator
}

func NewEventsUserHandler(bizkitEduApp *app.BizkitEduApp) *EventsUserHandler {
	eventsUserStorage := postgres.NewEventsUserStorage(bizkitEduApp.DB)
	eventsUserService := services.NewEventsUserService(eventsUserStorage)

	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	eventStorage := postgres.NewEventStorage(bizkitEduApp.DB)
	eventService := services.NewEventService(eventStorage)

	eventsUserUseCase := event.NewEventsUserUserCase(eventsUserService, userService, eventService)

	return &EventsUserHandler{
		eventsUserService: eventsUserService,
		eventsUserUseCase: eventsUserUseCase,
		validator:         bizkitEduApp.Validator,
	}
}

func (h EventsUserHandler) Create(c *fiber.Ctx) error {
	var eventsUser dto.CreateEventsUser

	if err := c.BodyParser(&eventsUser); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(eventsUser)
	if errValidate != nil {
		return errValidate
	}

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	eventsUserObj, err := h.eventsUserService.Create(c.Context(), &eventsUser, uuid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   eventsUserObj,
	})
}

func (h EventsUserHandler) GetAllByEventUUID(c *fiber.Ctx) error {
	var uuid4 apiDto.UUID
	uuid := c.Params("uuid")

	uuid4.UUID = uuid

	errValidate := h.validator.ValidateData(uuid4)
	if errValidate != nil {
		return errValidate
	}

	limit, offset := h.validator.GetLimitAndOffset(c)

	eventsUserObj, err := h.eventsUserUseCase.GetAllByEventUUID(c.Context(), uuid, limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   eventsUserObj,
	})
}

func (h EventsUserHandler) GetAllByUser(c *fiber.Ctx) error {
	var uuid4 apiDto.UUID
	uuid := c.Params("uuid")

	uuid4.UUID = uuid

	errValidate := h.validator.ValidateData(uuid4)
	if errValidate != nil {
		return errValidate
	}

	limit, offset := h.validator.GetLimitAndOffset(c)

	eventsUserObj, err := h.eventsUserUseCase.GetAllByUser(c.Context(), uuid, limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   eventsUserObj,
	})
}

func (h EventsUserHandler) Delete(c *fiber.Ctx) error {
	var uuid4 apiDto.UUID
	uuid := c.Params("uuid")

	uuid4.UUID = uuid

	errValidate := h.validator.ValidateData(uuid4)
	if errValidate != nil {
		return errValidate
	}

	err := h.eventsUserService.Delete(c.Context(), uuid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "successful delete event",
	})
}

func (h EventsUserHandler) Setup(router fiber.Router, handler fiber.Handler) {
	eventsUserGroup := router.Group("/eventsusers")
	eventsUserGroup.Use(handler)
	eventsUserGroup.Get("/getallbyevent/:uuid", h.GetAllByEventUUID, handler)
	eventsUserGroup.Post("/create", h.Create, handler)
	eventsUserGroup.Get("/getallbyuser/:uuid", h.GetAllByUser, handler)
}
