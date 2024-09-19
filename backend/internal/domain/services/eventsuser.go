package services

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/google/uuid"
)

type EventsUserStorage interface {
	Create(ctx context.Context, eventsUser *entities.EventsUser) (*entities.EventsUser, error)
	GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) ([]*entities.EventsUser, error)
	GetAllByUser(ctx context.Context, userUUID string, limit, offset int) ([]*entities.EventsUser, error)
	DeleteAllByEventUUID(ctx context.Context, eventUUID string) error
	DeleteByUserUUID(ctx context.Context, userUUID string) error
	Delete(ctx context.Context, UUID string) error
	IsRegistered(ctx context.Context, userUUID, eventUUID string) (*entities.EventsUser, error)
}

type eventsUserService struct {
	storage EventsUserStorage
}

func NewEventsUserService(storage EventsUserStorage) *eventsUserService {
	return &eventsUserService{storage: storage}
}

func (s *eventsUserService) Create(ctx context.Context, eventsUser *dto.CreateEventsUser, userUUID string) (*entities.EventsUser, error) {
	eventsUserObj := &entities.EventsUser{
		UUID:      uuid.NewString(),
		EventUUID: eventsUser.EventUUID,
		UserUUID:  userUUID,
	}

	return s.storage.Create(ctx, eventsUserObj)
}

func (s *eventsUserService) GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) ([]*entities.EventsUser, error) {
	return s.storage.GetAllByEventUUID(ctx, eventUUID, limit, offset)
}

func (s *eventsUserService) GetAllByUser(ctx context.Context, userUUID string, limit, offset int) ([]*entities.EventsUser, error) {
	return s.storage.GetAllByUser(ctx, userUUID, limit, offset)
}

func (s *eventsUserService) DeleteAllByEventUUID(ctx context.Context, eventUUID string) error {
	return s.storage.DeleteAllByEventUUID(ctx, eventUUID)
}

func (s *eventsUserService) DeleteByUserUUID(ctx context.Context, userUUID string) error {
	return s.storage.DeleteByUserUUID(ctx, userUUID)
}

func (s *eventsUserService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(ctx, UUID)
}

func (s *eventsUserService) IsRegistered(ctx context.Context, userUUID, eventUUID string) (*entities.EventsUser, error) {
	return s.storage.IsRegistered(ctx, userUUID, eventUUID)
}
