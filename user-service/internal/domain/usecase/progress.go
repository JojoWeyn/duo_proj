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
	// Собираем уникальные даты для указанного типа сущности
	dateMap := make(map[string]bool)
	for _, progress := range progresses {
		// Фильтруем по типу сущности
		if progress.EntityType == entityType {
			dateStr := progress.CompletedAt.Format("2006-01-02")
			dateMap[dateStr] = true
		}
	}

	// Получаем даты из dateMap
	var dates []time.Time
	for dateStr := range dateMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}

	log.Println("Date map: ", dateMap)

	// Добавляем текущий день, если на нем был прогресс
	currentDate := time.Now().Truncate(24 * time.Hour) // Получаем текущую дату без времени
	dateStr := currentDate.Format("2006-01-02")
	if _, exists := dateMap[dateStr]; !exists {
		// Добавляем только если еще нет этой даты в списке
		dates = append(dates, currentDate)
	}

	log.Println("Date map with current day: ", dates)

	// Сортируем даты
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Печатаем даты, чтобы проверить их порядок
	log.Println("Sorted dates: ", dates)

	// Рассчитываем текущий стрик
	currentStreak := 0
	for i := len(dates) - 1; i > 0; i-- {
		// Если день идет подряд с предыдущим
		if dates[i].Sub(dates[i-1]).Hours() == 24 {
			currentStreak++
		} else {
			// Прерывание стрика
			break
		}
	}

	// Печатаем текущий стрик
	log.Println("Current streak: ", currentStreak)

	return currentStreak, nil
}
