package usecase

import (
	"context"

	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, uuid uuid.UUID, login string) error {
	user := entity.NewUser(uuid, login)

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUseCase) GetUser(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	return uc.userRepo.FindByUUID(ctx, uuid)
}

func (uc *UserUseCase) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.GetAll(ctx, limit, offset)
}
func (uc *UserUseCase) UpdateUser(ctx context.Context, uuid uuid.UUID, updateData *entity.User) error {
	user, err := uc.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if updateData.Login != "" {
		user.Login = updateData.Login
	}

	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.SecondName != "" {
		user.SecondName = updateData.SecondName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.RankID != 0 {
		user.RankID = updateData.RankID
	}
	if updateData.Avatar != "" {
		user.Avatar = updateData.Avatar
	}

	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	return uc.userRepo.Delete(ctx, uuid)
}
