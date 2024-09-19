package postgres

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"gorm.io/gorm"
)

type eventsUserStorage struct {
	db *gorm.DB
}

func NewEventsUserStorage(db *gorm.DB) *eventsUserStorage {
	return &eventsUserStorage{db: db}
}

func (s *eventsUserStorage) Create(ctx context.Context, eventsUser *entities.EventsUser) (*entities.EventsUser, error) {
	err := s.db.WithContext(ctx).Create(&eventsUser).Error
	return eventsUser, err
}

func (s *eventsUserStorage) GetAllByEventUUID(ctx context.Context, eventUUID string, limit, offset int) ([]*entities.EventsUser, error) {
	var eventsUsers []*entities.EventsUser
	query := s.db.WithContext(ctx).Model(&entities.EventsUser{}).Order("created_at desc")

	err := query.Limit(limit).Offset(offset).Where("event_uuid = ?", eventUUID).Find(&eventsUsers).Error
	return eventsUsers, err
}

func (s *eventsUserStorage) GetAllByUser(ctx context.Context, userUUID string, limit, offset int) ([]*entities.EventsUser, error) {
	var eventsUsers []*entities.EventsUser
	query := s.db.WithContext(ctx).Model(&entities.EventsUser{}).Order("created_at desc")

	err := query.Limit(limit).Offset(offset).Where("user_uuid = ?", userUUID).Find(&eventsUsers).Error
	return eventsUsers, err
}

func (s *eventsUserStorage) DeleteAllByEventUUID(ctx context.Context, eventUUID string) error {
	err := s.db.WithContext(ctx).Unscoped().Delete(&entities.EventsUser{}, "event_uuid = ?", eventUUID).Error
	return err
}

func (s *eventsUserStorage) DeleteByUserUUID(ctx context.Context, userUUID string) error {
	err := s.db.WithContext(ctx).Unscoped().Delete(&entities.EventsUser{}, "user_id = ?", userUUID).Error
	return err
}

func (s *eventsUserStorage) Delete(ctx context.Context, UUID string) error {
	err := s.db.WithContext(ctx).Unscoped().Delete(&entities.EventsUser{}, "uuid = ?", UUID).Error
	return err
}

func (s *eventsUserStorage) IsRegistered(ctx context.Context, userUUID, eventUUID string) (*entities.EventsUser, error) {
	err := s.db.WithContext(ctx).Model(&entities.EventsUser{}).
		Where("user_uuid = ? AND event_uuid = ?", userUUID, eventUUID).
		Find(&entities.EventsUser{}).Error
	return &entities.EventsUser{}, err
}
