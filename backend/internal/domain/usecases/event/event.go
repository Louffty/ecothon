package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"io/ioutil"
	"net/http"
)

type Service interface {
	GetAll(ctx context.Context, limit, offset int, searchType string) ([]*entities.Event, error)
	GetAllByUserUUID(ctx context.Context, userUUID string, limit, offset int) ([]*entities.Event, error)
	GetUsersEvents(ctx context.Context, uuid string, limit, offset int) ([]*entities.Event, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Event, error)
}

type UserService interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
}

type EventsUserUseCase interface {
	IsRegistered(ctx context.Context, userUUID, eventUUID string) (bool, error)
}

type eventUseCase struct {
	eventService      Service
	userService       UserService
	eventsUserUseCase EventsUserUseCase
}

func NewEventUseCase(eventService Service, userService UserService, useCase EventsUserUseCase) *eventUseCase {
	return &eventUseCase{
		eventService:      eventService,
		userService:       userService,
		eventsUserUseCase: useCase,
	}
}

func (u eventUseCase) GetAll(ctx context.Context, limit, offset int, searchType, userUUID string) ([]*dto.ReturnEvent, error) {
	var (
		eventsDto []*dto.ReturnEvent
	)

	events, err := u.eventService.GetAll(ctx, limit, offset, searchType)

	if err != nil {
		return nil, err
	}

	for _, event := range events {
		user, errGetUser := u.userService.GetByUUID(ctx, event.AuthorUUID)

		if errGetUser != nil {
			return nil, errGetUser
		}

		isRegistered, err := u.eventsUserUseCase.IsRegistered(ctx, userUUID, event.UUID)
		if err != nil {
			return nil, err
		}

		eventsDto = append(eventsDto, &dto.ReturnEvent{
			Position: [2]float64{
				event.Latitude, event.Longitude,
			},
			Data: dto.Event{
				UUID:             event.UUID,
				Title:            event.Title,
				Description:      event.Description,
				StartTime:        event.StartTime,
				Longitude:        event.Longitude,
				Latitude:         event.Latitude,
				Address:          event.Address,
				IsUserRegistered: isRegistered,
				Author: dto.Author{
					UUID:       user.UUID,
					Username:   user.Username,
					Rate:       user.Rate,
					IsVerified: user.IsVerified,
				},
			},
		})
	}

	return eventsDto, nil
}

func (u eventUseCase) GetAllByUserUUID(ctx context.Context, userUUID string, limit, offset int) ([]*dto.Event, error) {
	var (
		eventsDto []*dto.Event
	)

	events, err := u.eventService.GetUsersEvents(ctx, userUUID, limit, offset)
	if err != nil {
		return nil, err
	}

	user, errGetUser := u.userService.GetByUUID(ctx, userUUID)
	if errGetUser != nil {
		return nil, errGetUser
	}

	for _, event := range events {
		eventsDto = append(eventsDto, &dto.Event{
			UUID:        event.UUID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime,
			Longitude:   event.Longitude,
			Latitude:    event.Latitude,
			Address:     event.Address,
			Author: dto.Author{
				UUID:       user.UUID,
				Username:   user.Username,
				Rate:       user.Rate,
				IsVerified: user.IsVerified,
			},
		})
	}

	return eventsDto, err
}

func (u eventUseCase) getRecommendationsFromExternalService(ctx context.Context, userUUID string, limit, offset int) ([]string, error) {
	type RecommendationRequest struct {
		UserUUID string `json:"user_uuid"`
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
	}

	// URL внешнего API рекомендаций
	apiURL := "https://nothypeproduction.space/recommendations/recommend"

	// Подготовка данных запроса
	reqBody := RecommendationRequest{
		UserUUID: userUUID,
		Limit:    limit,
		Offset:   offset,
	}

	// Конвертация структуры в JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Создание нового HTTP-запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")

	// Отправка HTTP-запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK HTTP status: %s. Response: %s", resp.Status, string(bodyBytes))
	}

	// Чтение тела ответа
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Распарсивание ответа в массив строк (UUID рекомендованных событий)
	var recommendations []string
	if err := json.Unmarshal(bodyBytes, &recommendations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	return recommendations, nil
}

func (u eventUseCase) GetRecommendation(ctx context.Context, userUUID string, limit, offset int) ([]*dto.Event, error) {
	var eventsDto []*dto.Event

	// Отправляем запрос к внешнему сервису рекомендаций
	recommendations, err := u.getRecommendationsFromExternalService(ctx, userUUID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, eventUUID := range recommendations {
		// Получаем событие по его UUID
		event, err := u.eventService.GetByUUID(ctx, eventUUID)
		if err != nil {
			return nil, err
		}

		// Получаем информацию об авторе события
		user, err := u.userService.GetByUUID(ctx, event.AuthorUUID)
		if err != nil {
			return nil, err
		}

		// Проверяем, зарегистрирован ли пользователь на это событие
		isRegistered, err := u.eventsUserUseCase.IsRegistered(ctx, userUUID, event.UUID)
		if err != nil {
			return nil, err
		}

		// Собираем DTO события
		eventsDto = append(eventsDto, &dto.Event{
			UUID:             event.UUID,
			Title:            event.Title,
			Description:      event.Description,
			StartTime:        event.StartTime,
			Longitude:        event.Longitude,
			Latitude:         event.Latitude,
			Address:          event.Address,
			IsUserRegistered: isRegistered,
			Author: dto.Author{
				UUID:       user.UUID,
				Username:   user.Username,
				Rate:       user.Rate,
				IsVerified: user.IsVerified,
			},
		})
	}

	return eventsDto, nil
}
