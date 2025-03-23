package usecase

import (
	"context"
	"errors"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type AttemptRepository interface {
	GetAttemptsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]entity.Attempt, error)
	CreateSessionAttempt(ctx context.Context, sessionAttempt *entity.AttemptSession) error
	CreateAttempt(ctx context.Context, attempt *entity.Attempt) error
	GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error)
	UpdateAttemptSession(ctx context.Context, sessionAttempt *entity.AttemptSession) error
	GetAttemptSessionByID(ctx context.Context, sessionUUID uuid.UUID) (*entity.AttemptSession, error)
}

type AttemptUseCase struct {
	repo AttemptRepository
}

func NewAttemptUseCase(repo AttemptRepository) *AttemptUseCase {
	return &AttemptUseCase{
		repo: repo,
	}
}

func (a *AttemptUseCase) StartAttempt(ctx context.Context, userUUID, exerciseUUID uuid.UUID) (*uuid.UUID, error) {
	sessionAttempt := entity.NewAttemptSession(userUUID, exerciseUUID)
	if err := a.repo.CreateSessionAttempt(ctx, sessionAttempt); err != nil {
		return nil, err
	}
	return &sessionAttempt.UUID, nil
}
func (a *AttemptUseCase) SubmitAnswer(ctx context.Context, sessionID, userUUID, questionUUID uuid.UUID, answer string, isCorrect bool) error {
	attempt := entity.NewAttempt(
		userUUID,
		sessionID,
		questionUUID,
		answer,
		isCorrect,
	)
	return a.repo.CreateAttempt(ctx, attempt)
}

func (a *AttemptUseCase) CreateAttempt(ctx context.Context, userUUID, questionUUID uuid.UUID, answer string, isCorrect bool) error {
	attempt := entity.NewAttempt(
		userUUID,
		uuid.UUID{},
		questionUUID,
		answer,
		isCorrect,
	)
	return a.repo.CreateAttempt(ctx, attempt)
}

func (a *AttemptUseCase) FinishAttempt(ctx context.Context, sessionID uuid.UUID) (int, int, error) {
	session, err := a.repo.GetAttemptSessionByID(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}

	if session.FinishedAt != nil {
		return 0, 0, errors.New("session already finished")
	}

	finishTime := time.Now()
	session.FinishedAt = &finishTime
	if err := a.repo.UpdateAttemptSession(ctx, session); err != nil {
		return 0, 0, err
	}

	attempts, err := a.repo.GetAttemptsBySessionID(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}

	correctAnswers := 0
	for _, attempt := range attempts {
		if attempt.IsCorrect {
			correctAnswers++
		}
	}
	questions := len(attempts)

	return correctAnswers, questions, nil

}

func (a *AttemptUseCase) GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error) {
	return a.repo.GetAttemptsByUser(ctx, userUUID)
}
