package services

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/google/uuid"
)

type AdminStorage interface {
	Create(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Admin, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Admin, error)
	Update(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
	Delete(ctx context.Context, uuid string) error
	GetByName(ctx context.Context, name string) (*entities.Admin, error)
}

type adminService struct {
	storage AdminStorage
}

func NewAdminService(storage AdminStorage) *adminService {
	return &adminService{storage: storage}
}

func (s *adminService) Create(ctx context.Context, createAdmin *dto.CreateAdminValue) (*entities.Admin, error) {
	admin := &entities.Admin{
		UUID:  uuid.NewString(),
		Name:  createAdmin.Name,
		Value: createAdmin.Value,
	}

	return s.storage.Create(ctx, admin)
}

func (s *adminService) GetByUUID(ctx context.Context, uuid string) (*entities.Admin, error) {
	return s.storage.GetByUUID(ctx, uuid)
}

func (s *adminService) Update(ctx context.Context, admin *dto.UpdateAdminValue) (*entities.Admin, error) {
	adminValue, err := s.storage.GetByUUID(ctx, admin.UUID)
	if err != nil {
		return nil, err
	}

	updatedValue := &entities.Admin{
		UUID:  admin.UUID,
		Name:  adminValue.Name,
		Value: admin.Value,
	}

	return s.storage.Update(ctx, updatedValue)
}

func (s *adminService) Delete(ctx context.Context, uuid string) error {
	return s.storage.Delete(ctx, uuid)
}

func (s *adminService) GetAll(ctx context.Context, limit, offset int) ([]*entities.Admin, error) {
	return s.storage.GetAll(ctx, limit, offset)
}

func (s *adminService) GetByName(ctx context.Context, name string) (*entities.Admin, error) {
	return s.storage.GetByName(ctx, name)
}
