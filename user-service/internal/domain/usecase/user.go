package usecase

import (
	"context"
	"fmt"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/kafka"
	"mime/multipart"

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

type StorageS3 interface {
	UploadAvatar(ctx context.Context, avatarFile multipart.File, fileName string, fileSize int64) (string, error)
}

type UserUseCase struct {
	userRepo UserRepository
	s3       StorageS3
	producer *kafka.Producer
}

func NewUserUseCase(userRepo UserRepository, s3 StorageS3, producer *kafka.Producer) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		s3:       s3,
		producer: producer,
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

	if err := user.Validate(); err != nil {
		return err
	}

	err = uc.producer.SendUserEvent(user.UUID.String(), user.Login, "update")
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return uc.userRepo.Update(ctx, user)
}
func (uc *UserUseCase) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	return uc.userRepo.Delete(ctx, uuid)
}

func (uc *UserUseCase) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarFile multipart.File, fileSize int64) (string, error) {
	fileName := fmt.Sprintf("%s_avatar.%s", userID, "png")

	avatarURL, err := uc.s3.UploadAvatar(ctx, avatarFile, fileName, fileSize)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	user, err := uc.userRepo.FindByUUID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	user.Avatar = avatarURL
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return "", fmt.Errorf("failed to update user avatar: %w", err)
	}

	return avatarURL, nil
}
