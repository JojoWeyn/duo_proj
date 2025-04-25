package postgres

import (
	"context"

	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Preload("Rank").Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}

	var totalPoints int
	entityType := "lesson"
	err := r.db.WithContext(ctx).
		Table("progresses").
		Select("COALESCE(SUM(points), 0)").
		Where("user_uuid = ? AND entity_type = ?", uuid, entityType).
		Scan(&totalPoints).Error
	if err != nil {
		return nil, err
	}
	user.TotalPoints = totalPoints

	var finished int64
	finishedType := "course"
	err = r.db.WithContext(ctx).
		Table("progresses").
		Where("user_uuid = ? AND entity_type = ? AND completed_at IS NOT NULL", uuid, finishedType).
		Count(&finished).Error
	if err != nil {
		return nil, err
	}
	user.FinishedCourses = finished

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.WithContext(ctx).
		Preload("Rank").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetLeaderboard(ctx context.Context, limit, offset int) ([]entity.Leaderboard, error) {
	var leaderboard []entity.Leaderboard

	query := `
		WITH user_points AS (
			SELECT
				u.uuid AS user_uuid,
				u.login,
				u.name,
				u.second_name,
				u.last_name,
				u.avatar,
				COALESCE(SUM(p.points), 0) AS total_points,
				MAX(p.completed_at) AS latest_completed_at -- используем для сортировки
			FROM
				users u
			LEFT JOIN
				progresses p ON u.uuid = p.user_uuid AND p.entity_type = 'lesson'
			GROUP BY
				u.uuid, u.login, u.name, u.second_name, u.last_name, u.avatar
		)
		SELECT
			up.user_uuid,
			up.login,
			up.name,
			up.second_name,
			up.last_name,
			up.avatar,
			up.total_points,
			ROW_NUMBER() OVER (ORDER BY up.total_points DESC, up.latest_completed_at DESC) AS rank
		FROM
			user_points up
		ORDER BY
			rank
		LIMIT ? OFFSET ?;`

	if err := r.db.WithContext(ctx).Raw(query, limit, offset).Scan(&leaderboard).Error; err != nil {
		return nil, err
	}

	return leaderboard, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("uuid = ?", user.UUID).
		Updates(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	return r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&entity.User{}).Error
}
