package setup

import (
	app "github.com/Louffty/green-code-moscow/cmd/app"
	v1 "github.com/Louffty/green-code-moscow/internal/adapters/controller/api/v1"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/v1/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// OAuth2 Configurations

func Setup(app *app.BizkitEduApp) {
	app.Fiber.Use(cors.New(cors.ConfigDefault))

	if app.Logging {
		app.Fiber.Use(logger.New())
	}

	app.Fiber.Get("/ping", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"body":   "pong",
		})
	})

	// Setup api v1 routes
	apiV1 := app.Fiber.Group("/api/v1")

	middlewareHandler := middlewares.NewMiddlewareHandler(app)

	// Setup user routes
	userHandler := v1.NewUserHandler(app)
	userHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

	// Setup question routes
	questionHandler := v1.NewQuestionHandler(app)
	questionHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

	// Setup conference routes
	conferenceHandler := v1.NewConferenceHandler(app)
	conferenceHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

	eventHandler := v1.NewEventHandler(app)
	eventHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

	adminHandler := v1.NewAdminHandler(app)
	adminHandler.Setup(apiV1, middlewareHandler.IsAdmin)

	eventsUser := v1.NewEventsUserHandler(app)
	eventsUser.Setup(apiV1, middlewareHandler.IsAuthenticated)

	oauthhandler := v1.NewOAuthHandler(app)
	oauthhandler.Setup(apiV1, middlewareHandler.IsAuthenticated)
}
