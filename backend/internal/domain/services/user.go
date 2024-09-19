package services

import (
	"context"
	"github.com/Louffty/green-code-moscow/internal/domain/common/errroz"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"github.com/Louffty/green-code-moscow/internal/domain/utils"
	"github.com/google/uuid"
)

// UserStorage is an interface that contains methods to interact with the database.
type UserStorage interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, uuid string) error
	GetByUsernameAndPassword(ctx context.Context, username string, password string) (*entities.User, error)
	Transfer(ctx context.Context, fromUUID, toUUID string, amount uint) error
	GetByName(ctx context.Context, username string) (*entities.User, error)
}

// userService is a struct that contains a pointer to an UserStorage instance.
type userService struct {
	adminService *adminService
	storage      UserStorage
}

// NewUserService is a function that returns a new instance of userService.
func NewUserService(storage UserStorage, adminStorage AdminStorage) *userService {
	return &userService{storage: storage, adminService: NewAdminService(adminStorage)}
}

func (s userService) Create(ctx context.Context, createUser *dto.CreateUser) (*entities.User, error) {
	user := &entities.User{
		UUID:     uuid.NewString(),
		Username: createUser.Username,
		Email:    createUser.Email,
	}
	user.SetPassword(createUser.Password)

	return s.storage.Create(ctx, user)
}

func (s userService) GenerateJwt(ctx context.Context, authUser *dto.AuthUser) (string, error) {
	user, err := s.storage.GetByUsernameAndPassword(ctx, authUser.Username, authUser.Password)
	if err != nil {
		return "", err
	}
	return utils.GenerateJwt(user.UUID, string(user.Password))
}

func (s userService) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	return s.storage.GetByUUID(ctx, uuid)
}

func (s userService) ChangeBalance(ctx context.Context, uuid string, change int) (*entities.User, error) {
	isAllFree, err := s.adminService.GetByName(ctx, "is_all_free")
	if err != nil {
		return &entities.User{}, err
	}

	if isAllFree.Value {
		return &entities.User{}, nil
	}

	user, err := s.storage.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if user.CoinsAmount+change < 0 {
		return nil, errroz.NotEnoughCoins
	}

	user.CoinsAmount += change
	return s.storage.Update(ctx, user)
}

// Transfer is a method to transfer coins between users.
func (s userService) Transfer(ctx context.Context, fromUUID, toUUID string, amount uint) error {
	if fromUUID == toUUID {
		return errroz.TransferToYourself
	}

	return s.storage.Transfer(ctx, fromUUID, toUUID, amount)
}

func (s userService) GetByName(ctx context.Context, username string) (*entities.User, error) {
	return s.storage.GetByName(ctx, username)
}

func (s userService) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	return s.storage.Update(ctx, user)
}
