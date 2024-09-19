package v1

import (
	"context"
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/validator"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/common/errroz"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/Louffty/green-code-moscow/internal/domain/usecases/event"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/gofiber/fiber/v2"
	"time"
)

type EventService interface {
	Create(ctx context.Context, event *dto.CreateEvent) (*entities.Event, error)
	GetAll(ctx context.Context, limit, offset int, searchType string) ([]*entities.Event, error)
	Archive(ctx context.Context, uuid string, userUIID string) (*entities.Event, error)
	GetUsersEvents(ctx context.Context, uuid string, limit, offset int) ([]*entities.Event, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Event, error)
}

type EventUseCase interface {
	GetAll(ctx context.Context, limit, offset int, searchType, userUUID string) ([]*dto.ReturnEvent, error)
	GetAllByUserUUID(ctx context.Context, userUUID string, limit, offset int) ([]*dto.Event, error)
	GetRecommendation(ctx context.Context, userUUID string, limit, offset int) ([]*dto.Event, error)
}

type EventHandler struct {
	eventService EventService
	eventUseCase EventUseCase
	validator    *validator.Validator
}

func NewEventHandler(bizkitEduApp *app.BizkitEduApp) *EventHandler {
	eventStorage := postgres.NewEventStorage(bizkitEduApp.DB)
	eventService := services.NewEventService(eventStorage)

	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	eventsUserStorage := postgres.NewEventsUserStorage(bizkitEduApp.DB)
	eventsUserService := services.NewEventsUserService(eventsUserStorage)
	eventsUserUseCase := event.NewEventsUserUserCase(eventsUserService, userService, eventService)

	eventUseCase := event.NewEventUseCase(eventService, userService, eventsUserUseCase)

	return &EventHandler{
		eventService: eventService,
		eventUseCase: eventUseCase,
		validator:    bizkitEduApp.Validator,
	}
}

func (h EventHandler) create(c *fiber.Ctx) error {
	var (
		createEvent dto.CreateEvent
		data        map[string]interface{}
	)

	if err := c.BodyParser(&createEvent); err != nil {
		return err
	}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	startTime, ok := data["time"]
	if !ok {
		return errroz.ParsingError
	}

	parsedTime, err := time.Parse("2006-01-02T15:04", startTime.(string))
	if err != nil {
		return err
	}
	createEvent.StartTime = parsedTime

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}
	createEvent.AuthorUUID = uuid

	if createEvent.StartTime.Before(time.Now()) {
		return errroz.InvalidStartTime
	}

	errValidate := h.validator.ValidateData(createEvent)
	if errValidate != nil {
		return errValidate
	}

	eventObject, err := h.eventService.Create(c.Context(), &createEvent)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   eventObject,
	})
}

func (h EventHandler) getAll(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	events, err := h.eventUseCase.GetAll(c.Context(), limit, offset, "upcoming", uuid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   events,
	})
}

func (h EventHandler) GetMyEvents(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	events, err := h.eventUseCase.GetAllByUserUUID(c.Context(), uuid, limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   events,
	})
}

func (h EventHandler) GetRecommend(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	events, err := h.eventUseCase.GetRecommendation(c.Context(), uuid, limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   events,
	})

}

func (h EventHandler) Setup(router fiber.Router, handler fiber.Handler) {
	eventGroup := router.Group("/event")
	eventGroup.Get("/all", h.getAll, handler)
	eventGroup.Post("/create", h.create, handler)
	eventGroup.Get("/allMy", h.GetMyEvents, handler)
	eventGroup.Get("/myRec", h.GetRecommend, handler)
}
