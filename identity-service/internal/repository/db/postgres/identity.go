package postgres

import (
	"context"
	"fmt"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"gorm.io/gorm"
)

type IdentityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) *IdentityRepository {
	return &IdentityRepository{
		db: db,
	}
}

func (r *IdentityRepository) Create(ctx context.Context, identity *entity.Identity) error {
	return r.db.WithContext(ctx).Create(identity).Error
}

func (r *IdentityRepository) FindByLogin(ctx context.Context, login string) (*entity.Identity, error) {
	var identity entity.Identity
	err := r.db.WithContext(ctx).Where("login = ?", login).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindByEmail(ctx context.Context, email string) (*entity.Identity, error) {
	var identity entity.Identity
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindByUUID(ctx context.Context, uuid string) (*entity.Identity, error) {
	var identity entity.Identity
	err := r.db.WithContext(ctx).Where("user_uuid = ?", uuid).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) Update(ctx context.Context, identity *entity.Identity) error {
	if identity.ID == 0 {
		return fmt.Errorf("cannot update entity: missing ID")
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Identity{}).
		Where("id = ?", identity.ID).
		Updates(identity).
		Error

	return err
}

func (r *IdentityRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.Identity{}, id).Error
}
