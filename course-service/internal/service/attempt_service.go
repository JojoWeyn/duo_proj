package service

import (
	"context"
	"errors"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"math"
	"time"
)

type QuestionRepo interface {
	GetCountByExerciseID(ctx context.Context, exerciseID uuid.UUID) (int64, error)
}

type ExerciseRepo interface {
	GetByQuestionID(ctx context.Context, questionID uuid.UUID) (*entity.Exercise, error)
}

type AttemptRepo interface {
	GetLast(ctx context.Context, userUUID, exerciseUUID uuid.UUID) ([]*entity.Attempt, error)
}

type AttemptService struct {
	questionRepo QuestionRepo
	exerciseRepo ExerciseRepo
	attemptRepo  AttemptRepo
}

func NewAttemptService(questionRepo QuestionRepo, exerciseRepo ExerciseRepo, attemptRepo AttemptRepo) *AttemptService {
	return &AttemptService{
		questionRepo: questionRepo,
		exerciseRepo: exerciseRepo,
		attemptRepo:  attemptRepo,
	}
}

func (a *AttemptService) ProcessUserAttempt(ctx context.Context, event kafka.UserAttemptEvent) (*kafka.UserProgressEvent, error) {
	exercise, err := a.exerciseRepo.GetByQuestionID(ctx, event.QuestionUUID)
	if err != nil {
		return nil, err
	}

	isCompleted, points := a.isExerciseCompleted(ctx, event.UserUUID, exercise.UUID)
	if isCompleted {
		progressEvent := kafka.UserProgressEvent{
			UserUUID:     event.UserUUID,
			ExerciseUUID: exercise.UUID,
			Points:       points,
			IsCorrect:    isCompleted,
			CreatedAt:    time.Now(),
		}
		return &progressEvent, nil
	}

	return nil, errors.New("exercise is not completed")
}

func (a *AttemptService) isExerciseCompleted(ctx context.Context, userUUID, exerciseUUID uuid.UUID) (bool, int) {
	totalQuestions, err := a.questionRepo.GetCountByExerciseID(ctx, exerciseUUID)
	if err != nil {
		return false, 0
	}

	lastAttempts, err := a.attemptRepo.GetLast(ctx, userUUID, exerciseUUID)
	if err != nil {
		return false, 0
	}

	correctCount := 0
	for _, attempt := range lastAttempts {
		if attempt.IsCorrect {
			correctCount++
		}
	}

	passThreshold := 1.0
	isPassed := float64(correctCount)/float64(totalQuestions) >= passThreshold

	points := int(math.Round(float64(correctCount) / float64(totalQuestions) * 100))

	return isPassed, points
}
