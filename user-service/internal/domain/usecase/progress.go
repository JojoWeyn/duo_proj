package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/google/uuid"
	"log"
	"sort"
	"time"
)

type ProgressRepository interface {
	GetProgress(ctx context.Context, userUUID uuid.UUID) ([]*entity.Progress, error)
	Create(ctx context.Context, progress *entity.Progress) error
	GetByEntityUUID(ctx context.Context, userUUID, entityUUID uuid.UUID) (*entity.Progress, error)
}
type ProgressUseCase struct {
	repo ProgressRepository
}

func NewProgressUseCase(repo ProgressRepository) *ProgressUseCase {
	return &ProgressUseCase{
		repo: repo,
	}
}

func (p *ProgressUseCase) GetProgress(ctx context.Context, userID uuid.UUID) ([]*entity.Progress, error) {
	var progresses []*entity.Progress

	progresses, err := p.repo.GetProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	return progresses, nil
}

func (p *ProgressUseCase) CheckProgress(ctx context.Context, userID, entityUUID uuid.UUID) bool {
	_, err := p.repo.GetByEntityUUID(ctx, userID, entityUUID)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (p *ProgressUseCase) AddProgress(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, points int, createdAt time.Time) error {
	progress := entity.NewProgress(userID, entityType, entityID, points, createdAt)

	if err := p.repo.Create(ctx, progress); err != nil {
		return err
	}

	return nil
}

func (p *ProgressUseCase) GetStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	progresses, err := p.repo.GetProgress(ctx, userID)
	if err != nil {
		return 0, err
	}

	entityType := "exercise"
	dateMap := make(map[string]bool)
	for _, progress := range progresses {
		if progress.EntityType == entityType {
			dateStr := progress.CompletedAt.Format("2006-01-02")
			dateMap[dateStr] = true
		}
	}

	var dates []time.Time
	for dateStr := range dateMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}

	// Сортируем по возрастанию
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	log.Println("Sorted dates: ", dates)

	// Стрик только если есть активность сегодня
	currentDate := time.Now().Truncate(24 * time.Hour)
	dateStr := currentDate.Format("2006-01-02")
	if !dateMap[dateStr] {
		log.Println("No progress today — streak is 0")
		return 0, nil
	}

	// Вычисляем стрик от сегодняшней даты назад
	currentStreak := 1
	for i := len(dates) - 2; i >= 0; i-- {
		// Сравниваем с предыдущей датой
		if currentDate.Sub(dates[i]).Hours() == float64(24*currentStreak) {
			currentStreak++
		} else {
			break
		}
	}

	log.Println("Current streak: ", currentStreak)
	return currentStreak, nil
}
