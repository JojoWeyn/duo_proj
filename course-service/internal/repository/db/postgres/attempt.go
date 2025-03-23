package postgres

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttemptRepository struct {
	db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) *AttemptRepository {
	return &AttemptRepository{
		db: db,
	}
}

func (a *AttemptRepository) CreateSessionAttempt(ctx context.Context, sessionAttempt *entity.AttemptSession) error {
	return a.db.WithContext(ctx).Create(sessionAttempt).Error
}

func (a *AttemptRepository) UpdateAttemptSession(ctx context.Context, sessionAttempt *entity.AttemptSession) error {
	return a.db.WithContext(ctx).Save(sessionAttempt).Error
}

func (a *AttemptRepository) GetAttemptSessionByID(ctx context.Context, sessionUUID uuid.UUID) (*entity.AttemptSession, error) {
	var attemptSession entity.AttemptSession
	if err := a.db.WithContext(ctx).
		Where("uuid = ?", sessionUUID).
		First(&attemptSession).Error; err != nil {
		return nil, err
	}
	return &attemptSession, nil
}

func (a *AttemptRepository) GetAttemptsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]entity.Attempt, error) {
	var attempts []entity.Attempt
	if err := a.db.WithContext(ctx).
		Where("attempt_session_uuid = ?", sessionID).
		Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (a *AttemptRepository) GetLast(ctx context.Context, userUUID, exerciseUUID, sessionUUID uuid.UUID) ([]*entity.Attempt, error) {
	var lastAttempts []*entity.Attempt
	if err := a.db.
		Raw(`
       WITH last_attempts AS (
			SELECT DISTINCT ON (question_uuid) * 
			FROM attempts 
			WHERE user_uuid = ? 
			AND question_uuid IN (SELECT uuid FROM questions WHERE exercise_uuid = ?)
			AND attempt_session_uuid = ?
			ORDER BY question_uuid, created_at DESC
		)
		SELECT * FROM last_attempts
	`, userUUID, exerciseUUID, sessionUUID).Scan(&lastAttempts).Error; err != nil {
		return nil, err
	}

	return lastAttempts, nil
}

func (a *AttemptRepository) CreateAttempt(ctx context.Context, attempt *entity.Attempt) error {
	return a.db.WithContext(ctx).Create(attempt).Error
}

func (a *AttemptRepository) GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error) {
	var attempts []entity.Attempt
	if err := a.db.WithContext(ctx).Where("user_uuid = ?", userUUID).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}
