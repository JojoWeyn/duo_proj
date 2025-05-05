package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

type AchievementRepository interface {
	GetAllAchievements(ctx context.Context) ([]entity.Achievement, error)
	GetAllUserAchievements(ctx context.Context, userID uuid.UUID) ([]dto.UserAchievementsDTO, error)
	GetUserAchievementProgress(ctx context.Context, userID uuid.UUID, achievementID int) (*entity.UserAchievementProgress, error)
	UpdateUserAchievementProgress(ctx context.Context, userID uuid.UUID, achievement entity.Achievement, countIncrement int) error
	CreateAchievement(ctx context.Context, achievement entity.Achievement) (entity.Achievement, error)
	GetAchievementByID(ctx context.Context, id int) (entity.Achievement, error)
	UpdateAchievement(ctx context.Context, id int, achievement entity.Achievement) (entity.Achievement, error)
	DeleteAchievement(ctx context.Context, id int) error
}

type AchievementUseCase struct {
	achievementRepo AchievementRepository
}

func NewAchievementUseCase(achievementRepo AchievementRepository) *AchievementUseCase {
	return &AchievementUseCase{
		achievementRepo: achievementRepo,
	}
}

func (uc *AchievementUseCase) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]dto.UserAchievementsDTO, error) {
	return uc.achievementRepo.GetAllUserAchievements(ctx, userID)
}

func (uc *AchievementUseCase) GetAllAchievements(ctx context.Context) ([]entity.Achievement, error) {
	return uc.achievementRepo.GetAllAchievements(ctx)
}

// CreateAchievement создает новое достижение
func (uc *AchievementUseCase) CreateAchievement(ctx context.Context, achievementData entity.Achievement) (entity.Achievement, error) {
	if err := achievementData.Validate(); err != nil {
		return entity.Achievement{}, err
	}
	
	return uc.achievementRepo.CreateAchievement(ctx, achievementData)
}

// GetAchievementByID получает достижение по ID
func (uc *AchievementUseCase) GetAchievementByID(ctx context.Context, id int) (entity.Achievement, error) {
	return uc.achievementRepo.GetAchievementByID(ctx, id)
}

func (uc *AchievementUseCase) UpdateAchievement(ctx context.Context, id int, achievementData entity.Achievement) (entity.Achievement, error) {
	// Проверяем существование достижения
	existingAchievement, err := uc.achievementRepo.GetAchievementByID(ctx, id)
	if err != nil {
		return entity.Achievement{}, err
	}

	// Обновляем поля
	achievementData.ID = existingAchievement.ID
	// Сохраняем дату создания из существующего достижения
	achievementData.CreatedAt = existingAchievement.CreatedAt

	// Обновляем достижение в репозитории
	return uc.achievementRepo.UpdateAchievement(ctx, id, achievementData)
}

func (uc *AchievementUseCase) DeleteAchievement(ctx context.Context, id int) error {
	// Проверяем существование достижения
	_, err := uc.achievementRepo.GetAchievementByID(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем достижение
	return uc.achievementRepo.DeleteAchievement(ctx, id)
}

func (uc *AchievementUseCase) CheckAchievements(ctx context.Context, userID uuid.UUID, action string) error {
	achievements, err := uc.achievementRepo.GetAllAchievements(ctx)
	if err != nil {
		return err
	}

	for _, ach := range achievements {
		var cond entity.Condition
		if err := json.Unmarshal([]byte(ach.Condition), &cond); err != nil {
			log.Printf("failed to unmarshal condition for achievement %d: %v", ach.ID, err)
			continue
		}

		if cond.Action != "" && cond.Action == action {
			uc.checkSimpleCounter(ctx, userID, ach, cond)
		} else if len(cond.ActionSeq) > 0 {
			uc.checkActionSequence(userID, ach, cond)
		} else if cond.Stat != "" {
			uc.checkStat(userID, ach, cond)
		}
	}
	return nil
}

func (uc *AchievementUseCase) checkSimpleCounter(ctx context.Context, userID uuid.UUID, ach entity.Achievement, cond entity.Condition) {
	countIncrement := 1

	err := uc.achievementRepo.UpdateUserAchievementProgress(ctx, userID, ach, countIncrement)
	if err != nil {
		fmt.Println("error updating user achievement progress:", err)
		return
	}

	progress, _ := uc.achievementRepo.GetUserAchievementProgress(ctx, userID, ach.ID)
	if progress != nil && progress.Achieved {
		fmt.Printf("user %s achieved: %s\n", userID, ach.Title)
	}
}

func (uc *AchievementUseCase) checkActionSequence(userID uuid.UUID, ach entity.Achievement, cond entity.Condition) {
	fmt.Println("check action sequence does not implemented yet")
}

func (uc *AchievementUseCase) checkStat(userID uuid.UUID, ach entity.Achievement, cond entity.Condition) {
	fmt.Println("check stat does not implemented yet")
}
