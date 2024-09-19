package question

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/common/errroz"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"sort"
)

// Service is an interface that contains a method to create a question.
type Service interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.Question, error)
	Create(ctx context.Context, question *dto.CreateQuestion) (*entities.Question, error)
	Close(ctx context.Context, uuid string) (*entities.Question, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Question, error)
	GetMy(ctx context.Context, limit, offset int, userUuid string) ([]*entities.Question, error)
}

// UserService is an interface that contains a method to change the balance of a user.
type UserService interface {
	ChangeBalance(ctx context.Context, uuid string, change int) (*entities.User, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
}

type AnswerService interface {
	GetAll(ctx context.Context, questionUUID string) ([]*entities.Answer, error)
	Create(ctx context.Context, answer *dto.CreateAnswer) (*entities.Answer, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Answer, error)
	Correct(ctx context.Context, answerUUID string) (*entities.Answer, error)
}

// question is a struct that contains instances of services.
type questionUseCase struct {
	questionService Service
	userService     UserService
	answerService   AnswerService
}

// NewQuestionUseCase is a function that returns a new instance of questionUseCase.
func NewQuestionUseCase(questionService Service, userService UserService, answerService AnswerService) *questionUseCase {
	return &questionUseCase{
		questionService: questionService,
		userService:     userService,
		answerService:   answerService,
	}
}

// CreateQuestion is a method that creates a new question.
func (u questionUseCase) CreateQuestion(ctx context.Context, question *dto.CreateQuestion) (*entities.Question, error) {
	_, err := u.userService.ChangeBalance(ctx, question.UserUUID, -int(question.Reward))
	if err != nil {
		return nil, err
	}

	return u.questionService.Create(ctx, question)
}

// CreateAnswer is a method that creates a new answer.
func (u questionUseCase) CreateAnswer(ctx context.Context, createAnswer *dto.CreateAnswer) (*entities.Answer, error) {
	_, err := u.userService.ChangeBalance(ctx, createAnswer.AuthorUUID, 10)

	question, err := u.questionService.GetByUUID(ctx, createAnswer.QuestionUUID)
	if err != nil {
		return nil, err
	}

	if question.Closed {
		return nil, errroz.QuestionClosed
	}

	return u.answerService.Create(ctx, createAnswer)
}

// GetQuestionWithAnswers is a method that returns a question with its answers.
func (u questionUseCase) GetQuestionWithAnswers(ctx context.Context, questionUUID string) (*entities.QuestionWithAnswers, error) {
	question, err := u.GetQuestionByUUID(ctx, questionUUID)
	if err != nil {
		return nil, err
	}

	answers, err := u.GetAllAnswersByUUID(ctx, questionUUID)
	if err != nil {
		return nil, err
	}

	return &entities.QuestionWithAnswers{
		Question: question,
		Answers:  answers,
	}, nil
}

// CorrectAnswer is a method for confirming the correctness of the response and closing the question with reward for answer author.
func (u questionUseCase) CorrectAnswer(ctx context.Context, answerUUID string, userUUID string) (*entities.QuestionWithAnswers, error) {
	answer, err := u.answerService.GetByUUID(ctx, answerUUID)
	if err != nil {
		return nil, err
	}

	question, err := u.questionService.GetByUUID(ctx, answer.QuestionUUID)
	if err != nil {
		return nil, err
	}

	if question.AuthorUUID != userUUID {
		return nil, errroz.NotEnoughPermissions
	}

	if question.Closed {
		return nil, errroz.QuestionClosed
	}

	question, err = u.questionService.Close(ctx, question.UUID)
	if err != nil {
		return nil, err
	}

	answer, err = u.answerService.Correct(ctx, answer.UUID)
	if err != nil {
		return nil, err
	}

	return u.GetQuestionWithAnswers(ctx, question.UUID)
}

func (u questionUseCase) GetAll(ctx context.Context, limit, offset int) ([]*dto.ReturnQuestion, error) {
	var (
		questionsWithFewAnswers  []*dto.ReturnQuestion
		questionsWithManyAnswers []*dto.ReturnQuestion
	)

	// Получаем все вопросы с пагинацией
	questions, err := u.questionService.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Итерируем по вопросам и собираем их в разные списки в зависимости от количества ответов
	for _, question := range questions {
		// Получаем количество ответов на вопрос
		answers, errGetAnswers := u.answerService.GetAll(ctx, question.UUID)
		if errGetAnswers != nil {
			return nil, err
		}

		// Получаем данные автора вопроса
		user, errGetUser := u.userService.GetByUUID(ctx, question.AuthorUUID)
		if errGetUser != nil {
			return nil, err
		}

		questionDto := &dto.ReturnQuestion{
			UUID:   question.UUID,
			Header: question.Header,
			Body:   question.Body,
			Reward: question.Reward,
			Author: dto.Author{
				UUID:     user.UUID,
				Username: user.Username,
				Rate:     user.Rate,
				Role:     user.Role,
			},
			Closed: question.Closed,
		}

		// Если ответов меньше 5, добавляем в первую группу
		if len(answers) < 5 {
			questionsWithFewAnswers = append(questionsWithFewAnswers, questionDto)
		} else {
			// Если ответов 5 или больше, добавляем во вторую группу
			questionsWithManyAnswers = append(questionsWithManyAnswers, questionDto)
		}
	}

	// Объединяем оба списка: сначала вопросы с < 5 ответами, потом с >= 5
	questionsDto := append(questionsWithFewAnswers, questionsWithManyAnswers...)

	return questionsDto, nil
}

func (u questionUseCase) GetAllAnswersByUUID(ctx context.Context, questionUUID string) ([]*dto.ReturnAnswer, error) {
	var (
		answerDto []*dto.ReturnAnswer
	)

	answers, err := u.answerService.GetAll(ctx, questionUUID)
	if err != nil {
		return nil, err
	}

	for _, answer := range answers {
		user, errGetUser := u.userService.GetByUUID(ctx, answer.AuthorUUID)
		if errGetUser != nil {
			return nil, err
		}

		answerDto = append(answerDto, &dto.ReturnAnswer{
			UUID:      answer.UUID,
			CreatedAt: answer.CreatedAt,
			UpdatedAt: answer.UpdatedAt,
			Author: dto.Author{
				UUID:     user.UUID,
				Username: user.Username,
				Rate:     user.Rate,
				Role:     user.Role,
			},
			QuestionUUID: answer.QuestionUUID,
			Body:         answer.Body,
			IsCorrect:    answer.IsCorrect,
		})
	}

	// Сортировка по роли, где сначала идут пользователи с ролью "admin"
	sort.SliceStable(answerDto, func(i, j int) bool {
		return answerDto[i].Author.Role == "admin" && answerDto[j].Author.Role != "admin"
	})

	return answerDto, nil
}

func (u questionUseCase) GetQuestionByUUID(ctx context.Context, uuid string) (*dto.ReturnQuestion, error) {
	question, err := u.questionService.GetByUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}

	user, errGetUser := u.userService.GetByUUID(ctx, question.AuthorUUID)

	if errGetUser != nil {
		return nil, err
	}

	return &dto.ReturnQuestion{
		UUID:   question.UUID,
		Header: question.Header,
		Body:   question.Body,
		Reward: question.Reward,
		Author: dto.Author{
			UUID:     user.UUID,
			Username: user.Username,
			Rate:     user.Rate,
			Role:     user.UUID,
		},
		Closed: question.Closed,
	}, nil
}

func (u questionUseCase) GetMyQuestions(ctx context.Context, limit int, offset int, uuid string) ([]*dto.ReturnQuestion, error) {
	var (
		questionsDto []*dto.ReturnQuestion
	)

	questions, err := u.questionService.GetMy(ctx, limit, offset, uuid)

	if err != nil {
		return nil, err
	}

	for _, question := range questions {
		user, errGetUser := u.userService.GetByUUID(ctx, question.AuthorUUID)

		if errGetUser != nil {
			return nil, errGetUser
		}

		questionsDto = append(questionsDto, &dto.ReturnQuestion{
			UUID:   question.UUID,
			Header: question.Header,
			Body:   question.Body,
			Reward: question.Reward,
			Author: dto.Author{
				UUID:     user.UUID,
				Username: user.Username,
				Rate:     user.Rate,
				Role:     user.Role,
			},
			Closed: question.Closed,
		})
	}

	return questionsDto, nil
}
