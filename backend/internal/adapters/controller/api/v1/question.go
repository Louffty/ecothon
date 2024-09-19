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
	"github.com/Louffty/green-code-moscow/internal/domain/usecases/question"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/gofiber/fiber/v2"
)

// QuestionService is an interface that contains methods to interact with the question service
type QuestionService interface {
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Question, error)
	GetMy(ctx context.Context, limit, offset int, userUuid string) ([]*entities.Question, error)
}

type QuestionUseCase interface {
	CreateQuestion(ctx context.Context, question *dto.CreateQuestion) (*entities.Question, error)
	CreateAnswer(ctx context.Context, createAnswer *dto.CreateAnswer) (*entities.Answer, error)
	GetQuestionWithAnswers(ctx context.Context, questionUUID string) (*entities.QuestionWithAnswers, error)
	CorrectAnswer(ctx context.Context, answerUUID string, userUUID string) (*entities.QuestionWithAnswers, error)
	GetAll(ctx context.Context, limit, offset int) ([]*dto.ReturnQuestion, error)
	GetMyQuestions(ctx context.Context, limit int, offset int, uuid string) ([]*dto.ReturnQuestion, error)
}

// QuestionHandler is a struct that contains the questionService and validator.
type QuestionHandler struct {
	questionService QuestionService
	questionUseCase QuestionUseCase
	validator       *validator.Validator
}

// NewQuestionHandler is a function that returns a new instance of QuestionHandler.
func NewQuestionHandler(bizkitEduApp *app.BizkitEduApp) *QuestionHandler {
	questionStorage := postgres.NewQuestionStorage(bizkitEduApp.DB)
	questionService := services.NewQuestionService(questionStorage)

	userStorage := postgres.NewUserStorage(bizkitEduApp.DB)
	adminStorage := postgres.NewAdminStorage(bizkitEduApp.DB)
	userService := services.NewUserService(userStorage, adminStorage)

	answerStorage := postgres.NewAnswerStorage(bizkitEduApp.DB)
	answerService := services.NewAnswerService(answerStorage)

	questionUseCase := question.NewQuestionUseCase(questionService, userService, answerService)

	return &QuestionHandler{
		questionService: questionService,
		questionUseCase: questionUseCase,
		validator:       bizkitEduApp.Validator,
	}
}

// create is a handler for creating new question in database.
func (h QuestionHandler) create(c *fiber.Ctx) error {
	var createQuestion dto.CreateQuestion

	if err := c.BodyParser(&createQuestion); err != nil {
		return err
	}

	errValidate := h.validator.ValidateData(createQuestion)
	if errValidate != nil {
		return errValidate
	}

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}
	createQuestion.UserUUID = uuid

	q, err := h.questionUseCase.CreateQuestion(c.Context(), &createQuestion)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   q,
	})
}

// getAll is a handler for getting all questions.
func (h QuestionHandler) getAll(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	questions, err := h.questionUseCase.GetAll(c.Context(), limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   questions,
	})
}

// getByUUID is a handler for getting a question by UUID.
func (h QuestionHandler) getByUUID(c *fiber.Ctx) error {
	var uuid4 apiDto.UUID
	uuid := c.Params("uuid")

	uuid4.UUID = uuid

	errValidate := h.validator.ValidateData(uuid4)
	if errValidate != nil {
		return errValidate
	}

	q, err := h.questionUseCase.GetQuestionWithAnswers(c.Context(), uuid4.UUID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   q,
	})
}

// createAnswer is a handler for creating an answer.
func (h QuestionHandler) createAnswer(c *fiber.Ctx) error {
	var createAnswer dto.CreateAnswer

	if err := c.BodyParser(&createAnswer); err != nil {
		return err
	}

	authorUuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}
	createAnswer.AuthorUUID = authorUuid

	answer, err := h.questionUseCase.CreateAnswer(c.Context(), &createAnswer)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body":   answer,
	})
}

// correctAnswer is a handler for confirming the correctness of the response.
func (h QuestionHandler) correctAnswer(c *fiber.Ctx) error {
	var uuid4 apiDto.UUID
	uuid := c.Params("uuid")

	uuid4.UUID = uuid

	errValidate := h.validator.ValidateData(uuid4)
	if errValidate != nil {
		return errValidate
	}

	userUUID, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	questionWithAnswers, err := h.questionUseCase.CorrectAnswer(c.Context(), uuid4.UUID, userUUID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   questionWithAnswers,
	})
}

// getMy is a handler for getting all questions of the user.
func (h QuestionHandler) getMy(c *fiber.Ctx) error {
	limit, offset := h.validator.GetLimitAndOffset(c)

	uuid, err := utils.GetUUIDByToken(c)
	if err != nil {
		return err
	}

	questions, err := h.questionUseCase.GetMyQuestions(c.Context(), limit, offset, uuid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body":   questions,
	})
}

// Setup is a function that registers all routes for the question.
func (h QuestionHandler) Setup(router fiber.Router, handler fiber.Handler) {
	questionGroup := router.Group("/question")
	questionGroup.Post("/create", h.create, handler)
	questionGroup.Get("/all", h.getAll, handler)
	questionGroup.Get("/my", h.getMy, handler)
	questionGroup.Get("/:uuid", h.getByUUID, handler)
	questionGroup.Post("/answer/create", h.createAnswer, handler)
	questionGroup.Put("/answer/correct/:uuid", h.correctAnswer, handler)
}
