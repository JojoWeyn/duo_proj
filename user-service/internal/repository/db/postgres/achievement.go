package postgres

import (
	"context"
	"encoding/json"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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
