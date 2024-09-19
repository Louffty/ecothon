package postgres

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"gorm.io/gorm"
)

// adminStorage is a struct that contains a pointer to a gorm.DB instance.
type adminStorage struct {
	db *gorm.DB
}

// NewAdminStorage is a function that returns a new instance of adminsStorage.
func NewAdminStorage(db *gorm.DB) *adminStorage {
	return &adminStorage{db: db}
}

// Create is a method to create a new User in database.
func (s *adminStorage) Create(ctx context.Context, admin *entities.Admin) (*entities.Admin, error) {
	err := s.db.WithContext(ctx).Create(&admin).Error
	return admin, err
}

// GetByUUID is a method that returns an error and a pointer to a User instance.
func (s *adminStorage) GetByUUID(ctx context.Context, uuid string) (*entities.Admin, error) {
	var admin *entities.Admin
	err := s.db.WithContext(ctx).Model(&entities.Admin{}).Where("uuid = ?", uuid).First(&admin).Error
	return admin, err
}

// GetAll is a method that returns a slice of pointers to User instances.
func (s *adminStorage) GetAll(ctx context.Context, limit, offset int) ([]*entities.Admin, error) {
	var admins []*entities.Admin
	err := s.db.WithContext(ctx).Model(&entities.Admin{}).Limit(limit).Offset(offset).Find(&admins).Error
	return admins, err
}

// Update is a method to update an existing User in database.
func (s *adminStorage) Update(ctx context.Context, admin *entities.Admin) (*entities.Admin, error) {
	err := s.db.WithContext(ctx).Model(&entities.Admin{}).Where("uuid = ?", admin.UUID).Updates(&admin).Error
	return admin, err
}

// Delete is a method to delete an existing User in database.
func (s *adminStorage) Delete(ctx context.Context, uuid string) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&entities.Admin{}, "uuid = ?", uuid).Error
}

func (s *adminStorage) GetByName(ctx context.Context, name string) (*entities.Admin, error) {
	var admin *entities.Admin
	err := s.db.WithContext(ctx).Model(&entities.Admin{}).Where("name = ?", name).First(&admin).Error
	return admin, err
}
