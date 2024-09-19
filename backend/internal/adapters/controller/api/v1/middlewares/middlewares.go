package middlewares

import (
	"context"
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/common/errroz"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type UserService interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
}

type MiddlewareHandler struct {
	userService UserService
}

// NewMiddlewareHandler is a function that returns a new instance of MiddlewareHandler.
func NewMiddlewareHandler(bizkitEduApp *app.BizkitEduApp) *MiddlewareHandler {
	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	return &MiddlewareHandler{
		userService: userService,
	}
}

func (h MiddlewareHandler) IsAuthenticated(c *fiber.Ctx) error {
	// Проверяем заголовок Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Если Authorization заголовок пуст, проверяем куки
		authCookie := c.Cookies("jwt_token")
		if authCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "auth header and cookie are empty",
			})
		}

		// Если токен найден в куках, используем его
		authHeader = "Bearer " + authCookie
	}

	// Разбираем заголовок Authorization
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "invalid auth header",
		})
	}

	// Извлекаем данные из JWT
	uuid, password, errParse := utils.ParseJwt(parts[1])
	if errParse != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errParse.Error(),
		})
	}

	// Ищем пользователя по UUID
	user, errGetUser := h.userService.GetByUUID(c.Context(), uuid)
	if errGetUser != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errGetUser.Error(),
		})
	}

	// Проверяем пароль пользователя
	if string(user.Password) != password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errroz.TokenExpired.Error(),
		})
	}

	// Если все проверки пройдены, продолжаем выполнение
	return c.Next()
}

func (h MiddlewareHandler) IsAdmin(c *fiber.Ctx) error {
	// Проверяем заголовок Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Если Authorization заголовок пуст, проверяем куки
		authCookie := c.Cookies("jwt_token")
		if authCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "auth header and cookie are empty",
			})
		}

		// Если токен найден в куках, используем его
		authHeader = "Bearer " + authCookie
	}

	// Разбираем заголовок Authorization
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "invalid auth header",
		})
	}

	// Извлекаем данные из JWT
	uuid, password, errParse := utils.ParseJwt(parts[1])
	if errParse != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errParse.Error(),
		})
	}

	// Ищем пользователя по UUID
	user, errGetUser := h.userService.GetByUUID(c.Context(), uuid)
	if errGetUser != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errGetUser.Error(),
		})
	}

	// Проверяем пароль пользователя
	if string(user.Password) != password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errroz.TokenExpired.Error(),
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "forbidden",
		})
	}

	return c.Next()
}

func (h MiddlewareHandler) IsMasterAdmin(c *fiber.Ctx) error {
	if len(c.GetReqHeaders()["Authorization"]) > 0 {
		authHeader := c.GetReqHeaders()["Authorization"][0]
		if authHeader == "" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "auth header is empty",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "invalid auth header",
			})
		}

		uuid, password, errParse := utils.ParseJwt(parts[1])
		if errParse != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"body":   errParse.Error(),
			})
		}

		user, errGetUser := h.userService.GetByUUID(c.Context(), uuid)
		if errGetUser != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"body":   errGetUser.Error(),
			})
		}

		if string(user.Password) != password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"body":   errroz.TokenExpired.Error(),
			})
		}

		if user.Role != "master_admin" {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "forbidden",
			})
		}

		return c.Next()
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": false,
		"body":   errroz.EmptyAuthHeader.Error(),
	})
}
