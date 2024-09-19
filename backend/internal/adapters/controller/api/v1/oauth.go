package v1

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Louffty/green-code-moscow/cmd/app"
	"github.com/Louffty/green-code-moscow/internal/adapters/config"
	"github.com/Louffty/green-code-moscow/internal/adapters/controller/api/validator"
	"github.com/Louffty/green-code-moscow/internal/adapters/database/postgres"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
)

type OAuthHandler struct {
	userService UserService
	validator   *validator.Validator
}

func NewOAuthHandler(bizkitEduApp *app.BizkitEduApp) *OAuthHandler {
	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	return &OAuthHandler{
		userService: userService,
		validator:   bizkitEduApp.Validator,
	}

}

func generateRandomState() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "randomstate" // fallback in case of error
	}
	return hex.EncodeToString(bytes)
}

func (h OAuthHandler) VKLogin(c *fiber.Ctx) error {
	conf := config.VKConfig()

	// Generate a random state
	state := generateRandomState()

	// Save state in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Path:     "/",
	})

	url := conf.AuthCodeURL(state)

	// Redirect user to VK OAuth
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)

	return nil
}

func (h OAuthHandler) VKCallback(c *fiber.Ctx) error {
	// Check state match
	state := c.Query("state")
	oauthStateCookie := c.Cookies("oauth_state")
	if state != oauthStateCookie {
		return c.Status(fiber.StatusBadRequest).SendString("OAuth state mismatch")
	}

	// Get code
	code := c.Query("code")

	// Exchange code for token
	vkConfig := config.VKConfig()
	token, err := vkConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Code-Token Exchange Failed: " + err.Error())
	}

	// Fetch user data from VK
	resp, err := http.Get("https://api.vk.com/method/users.get?access_token=" + token.AccessToken + "&v=5.131")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("User Data Fetch Failed: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to fetch user data: " + resp.Status)
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read user data: " + err.Error())
	}

	// Structure for VK user data
	type VKUserResponse struct {
		Response []struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"response"`
	}

	var vkUserResponse VKUserResponse
	err = json.Unmarshal(userData, &vkUserResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("User Data Unmarshal Failed: " + err.Error())
	}

	if len(vkUserResponse.Response) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No user data found")
	}

	vkUser := vkUserResponse.Response[0]

	if vkUser.Email == "" {
		vkUser.Email = strconv.Itoa(vkUser.ID)
	}

	pass := generateRandomState()

	// Check if user exists in the database
	user, err := h.userService.GetByName(c.Context(), strconv.Itoa(vkUser.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If user is not found, create them
			newUser := &dto.CreateUser{
				Username: strconv.Itoa(vkUser.ID),
				Email:    vkUser.Email,
				Password: pass, // You may want to handle passwords differently for VK users
			}

			user, err = h.userService.Create(c.Context(), newUser)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("User Registration Failed: " + err.Error())
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user: " + err.Error())
		}
	}

	user.SetPassword(pass)

	user, errUpdate := h.userService.Update(c.Context(), user)
	if errUpdate != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("User Update Failed: " + errUpdate.Error())
	}

	// Generate JWT
	authUser := dto.AuthUser{
		Username: user.Username,
		Password: pass, // Assuming no password for VK users
	}

	jwt, err := h.userService.GenerateJwt(c.Context(), &authUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("JWT Generation Failed: " + err.Error())
	}

	// Set JWT in cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "Bearer " + jwt,
		Path:     "/",
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// Redirect to external domain
	return c.Redirect("https://nothypeproduction.space")
}

func (h OAuthHandler) Setup(router fiber.Router, handler fiber.Handler) {
	oauthGroup := router.Group("/oauth")
	oauthGroup.Get("/vk_login", h.VKLogin)
	oauthGroup.Get("/vk_callback", h.VKCallback)
}
