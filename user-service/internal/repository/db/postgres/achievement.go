package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) GetAllAchievements(ctx context.Context) ([]entity.Achievement, error) {
	var achievements []entity.Achievement
	err := r.db.WithContext(ctx).Find(&achievements).Error
	return achievements, err
}

func (r *AchievementRepository) GetUserAchievementProgress(ctx context.Context, userID uuid.UUID, achievementID int) (*entity.UserAchievementProgress, error) {
	var progress entity.UserAchievementProgress
	err := r.db.WithContext(ctx).
		Where("user_uuid = ? AND achievement_id = ?", userID, achievementID).
		First(&progress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &progress, nil
}

func (r *AchievementRepository) UpdateUserAchievementProgress(ctx context.Context, userID uuid.UUID, achievement entity.Achievement, countIncrement int) error {
	now := time.Now()
	var progress entity.UserAchievementProgress

	err := r.db.WithContext(ctx).
		Where("user_uuid = ? AND achievement_id = ?", userID, achievement.ID).
		First(&progress).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			progress = entity.UserAchievementProgress{
				UserUUID:      userID,
				AchievementID: achievement.ID,
				CurrentCount:  countIncrement,
				Achieved:      false,
				UpdatedAt:     now,
			}
			return r.db.WithContext(ctx).Create(&progress).Error
		}
		return err
	}

	if progress.Achieved {
		return nil
	}

	progress.CurrentCount += countIncrement
	progress.UpdatedAt = now

	var cond entity.Condition
	err = json.Unmarshal([]byte(achievement.Condition), &cond)
	if err != nil {
		return err
	}

	if progress.CurrentCount >= cond.Count && !progress.Achieved {
		progress.Achieved = true
		progress.AchievedAt = &now
	}

	return r.db.WithContext(ctx).Save(&progress).Error
}

func (r *AchievementRepository) GetAllUserAchievements(ctx context.Context, userID uuid.UUID) ([]dto.UserAchievementsDTO, error) {
	var achievements []dto.UserAchievementsDTO

	err := r.db.Table("achievements AS a").
		WithContext(ctx).
		Select("a.id, a.title, a.description, a.secret, a.created_at, ua.achieved_at, ua.current_count, ua.achieved").
		Joins("JOIN user_achievement_progresses ua ON ua.achievement_id = a.id").
		Where("ua.user_uuid = ?", userID).
		Scan(&achievements).Error

	if err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementRepository) CreateAchievement(ctx context.Context, achievement entity.Achievement) (entity.Achievement, error) {
	achievement.ID = 0
	achievement.CreatedAt = time.Now()
	err := r.db.WithContext(ctx).Create(&achievement).Error
	return achievement, err
}

func (r *AchievementRepository) GetAchievementByID(ctx context.Context, id int) (entity.Achievement, error) {
	var achievement entity.Achievement
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&achievement).Error
	return achievement, err
}

func (r *AchievementRepository) UpdateAchievement(ctx context.Context, id int, achievement entity.Achievement) (entity.Achievement, error) {
	var existingAchievement entity.Achievement
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&existingAchievement).Error; err != nil {
		return entity.Achievement{}, err
	}

	existingAchievement.Title = achievement.Title
	existingAchievement.Description = achievement.Description
	existingAchievement.Condition = achievement.Condition
	existingAchievement.Secret = achievement.Secret

	err := r.db.WithContext(ctx).Save(&existingAchievement).Error
	return existingAchievement, err
}

func (r *AchievementRepository) DeleteAchievement(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&entity.Achievement{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
