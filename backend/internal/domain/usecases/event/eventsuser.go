package event

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
)

type EventsUserService interface {
	GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) ([]*entities.EventsUser, error)
	GetAllByUser(ctx context.Context, userUUID string, limit, offset int) ([]*entities.EventsUser, error)
	IsRegistered(ctx context.Context, userUUID, eventUUID string) (*entities.EventsUser, error)
}

type EventService interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.Event, error)
}

type eventsUserUseCase struct {
	eventsUserService EventsUserService
	eventService      EventService
	userService       UserService
}

func NewEventsUserUserCase(eventsUserService EventsUserService, userService UserService, eventService EventService) *eventsUserUseCase {
	return &eventsUserUseCase{
		eventsUserService: eventsUserService,
		eventService:      eventService,
		userService:       userService,
	}
}

func (u eventsUserUseCase) GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) (*dto.ReturnEventsUser, error) {
	var returnEventsUsers dto.ReturnEventsUser

	eventsUser, err := u.eventsUserService.GetAllByEventUUID(ctx, eventUUID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, eventUser := range eventsUser {
		user, errGetUser := u.userService.GetByUUID(ctx, eventUser.UserUUID)
		if errGetUser != nil {
			return nil, errGetUser
		}

		returnEventsUsers.Users = append(returnEventsUsers.Users, &dto.Author{
			UUID:     user.UUID,
			Username: user.Username,
			Rate:     user.Rate,
		})
	}

	event, errGetEvent := u.eventService.GetByUUID(ctx, eventUUID)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	author, err := u.userService.GetByUUID(ctx, event.AuthorUUID)
	if err != nil {
		return nil, err
	}

	returnEventsUsers.Event = dto.Event{
		UUID:        event.UUID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		Longitude:   event.Longitude,
		Latitude:    event.Latitude,
		Address:     event.Address,
		Author: dto.Author{
			UUID:     author.UUID,
			Username: author.Username,
			Rate:     author.Rate,
		},
	}

	return &returnEventsUsers, err
}

func (u eventsUserUseCase) GetAllByUser(ctx context.Context, userUUID string, limit, offset int) (*dto.ReturnUsersEvents, error) {
	var returnEventsUsers dto.ReturnUsersEvents

	eventsUser, err := u.eventsUserService.GetAllByUser(ctx, userUUID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, eventUser := range eventsUser {
		event, errGetEvent := u.eventService.GetByUUID(ctx, eventUser.EventUUID)
		if errGetEvent != nil {
			return nil, errGetEvent
		}

		author, errGetAuthor := u.userService.GetByUUID(ctx, event.AuthorUUID)
		if errGetAuthor != nil {
			return nil, errGetAuthor
		}

		returnEventsUsers.Events = append(returnEventsUsers.Events, &dto.Event{
			UUID:        event.UUID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime,
			Longitude:   event.Longitude,
			Latitude:    event.Latitude,
			Address:     event.Address,
			Author: dto.Author{
				UUID:     author.UUID,
				Username: author.Username,
				Rate:     author.Rate,
			},
		})
	}

	return &returnEventsUsers, err
}

func (u eventsUserUseCase) IsRegistered(ctx context.Context, userUUID, eventUUID string) (bool, error) {
	registered, err := u.eventsUserService.IsRegistered(ctx, userUUID, eventUUID)
	if err != nil {
		return false, err
	}

	// Проверка: если `registered` не nil и есть данные в каком-то поле, например `UserUUID`
	if registered != nil && registered.UserUUID != "" {
		return true, nil
	}

	// Если `registered` nil или поля пустые, значит регистрация отсутствует
	return false, nil
}
