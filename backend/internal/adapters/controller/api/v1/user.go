package v1

import (
	"context"
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/validator"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/Louffty/green-code-moscow/internal/domain/usecases/event"
	"github.com/Louffty/green-code-moscow/internal/domain/usecases/user"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/gofiber/fiber/v2"
)

// UserService is an interface that contains methods to interact with the user service
type UserService interface {
	Create(ctx context.Context, createUser *dto.CreateUser) (*entities.User, error)
	GenerateJwt(ctx context.Context, authUser *dto.AuthUser) (string, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
	GetByName(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
}

type UserUseCase interface {
	GetByUUID(ctx context.Context, uuid string) (*dto.ReturnUser, error)
	VerifiedUser(ctx context.Context, uuid string, organisation string) (*entities.User, error)
}

// UserHandler is a struct that contains the userService and validator.
type UserHandler struct {
	userService UserService
	userUseCase UserUseCase
	validator   *validator.Validator
}

// NewUserHandler is a function that returns a new instance of UserHandler.
func NewUserHandler(bizkitEduApp *app.BizkitEduApp) *UserHandler {
	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	eventsUserStorage := postgres.NewEventsUserStorage(bizkitEduApp.DB)
	eventsUserService := services.NewEventsUserService(eventsUserStorage)

	eventStorage := postgres.NewEventStorage(bizkitEduApp.DB)
	eventService := services.NewEventService(eventStorage)

	eventsUserUseCase := event.NewEventsUserUserCase(eventsUserService, userService, eventService)

	userUseCase := user.NewUserUseCase(userService, eventsUserUseCase)

	return &UserHandler{
		userService: userService,
		userUseCase: userUseCase,
		validator:   bizkitEduApp.Validator,
	}
}

// Register is handler for user registration.
func (h UserHandler) register(c *fiber.Ctx) error {
	var createUser dto.CreateUser

	if err := c.BodyParser(&createUser); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(createUser)
	if errValidate != nil {
		return errValidate
	}

	user, err := h.userService.Create(c.Context(), &createUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   user,
	})
}

func (h UserHandler) auth(c *fiber.Ctx) error {
	var authUser dto.AuthUser

	if err := c.BodyParser(&authUser); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(authUser)
	if errValidate != nil {
		return errValidate
	}

	jwt, err := h.userService.GenerateJwt(c.Context(), &authUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   jwt,
	})
}

func (h UserHandler) me(c *fiber.Ctx) error {
	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetByUUID(c.Context(), uuid)
	if err != nil {
		return err
	}

	returnUser, err := h.userUseCase.GetByUUID(c.Context(), user.UUID)

	return c.JSON(fiber.Map{
		"status": true,
		"body":   returnUser,
	})

}

func (h UserHandler) verified(c *fiber.Ctx) error {
	var org dto.VerifiedUser

	if err := c.BodyParser(&org); err != nil {
		return err
	}

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	verified, err := h.userUseCase.VerifiedUser(c.Context(), uuid, org.Organisation)
	if err != nil {
		return err
	}

	var stat bool

	if verified.IsVerified {
		stat = true
	} else {
		stat = false
	}

	return c.JSON(fiber.Map{
		"status": stat,
		"body":   verified,
	})
}

// Setup is a function that registers all routes for the user.
func (h UserHandler) Setup(router fiber.Router, handler fiber.Handler) {
	userGroup := router.Group("/user")
	userGroup.Post("/register", h.register)
	userGroup.Post("/auth", h.auth)
	userGroup.Get("/me", h.me, handler)
	userGroup.Post("/verified", h.verified, handler)
}
