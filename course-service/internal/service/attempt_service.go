package service

import (
	"context"
	"errors"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type CompletionRepo interface {
	MarkExerciseCompleted(ctx context.Context, userUUID, exerciseUUID uuid.UUID, points int) error
	IsExerciseCompleted(ctx context.Context, userUUID, exerciseUUID uuid.UUID) (bool, int, error)

	MarkLessonCompleted(ctx context.Context, userUUID, lessonUUID uuid.UUID, points int) error
	IsLessonCompleted(ctx context.Context, userUUID, lessonUUID uuid.UUID) (bool, int, error)

	MarkCourseCompleted(ctx context.Context, userUUID, courseUUID uuid.UUID, points int) error
	IsCourseCompleted(ctx context.Context, userUUID, courseUUID uuid.UUID) (bool, int, error)
}

type QuestionRepo interface {
	GetCountByExerciseID(ctx context.Context, exerciseID uuid.UUID) (int64, error)
}

type ExerciseRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error)
	GetByQuestionID(ctx context.Context, questionID uuid.UUID) (*entity.Exercise, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error)
}

type AttemptRepo interface {
	GetLast(ctx context.Context, userUUID, exerciseUUID, sessionUUID uuid.UUID) ([]*entity.Attempt, error)
}

type LessonRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error)
}

type AttemptService struct {
	questionRepo   QuestionRepo
	exerciseRepo   ExerciseRepo
	attemptRepo    AttemptRepo
	lessonRepo     LessonRepo
	completionRepo CompletionRepo
}

func NewAttemptService(questionRepo QuestionRepo, exerciseRepo ExerciseRepo, attemptRepo AttemptRepo, lessonRepo LessonRepo, completionRepo CompletionRepo) *AttemptService {
	return &AttemptService{
		questionRepo:   questionRepo,
		exerciseRepo:   exerciseRepo,
		attemptRepo:    attemptRepo,
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
	}
}

func (a *AttemptService) ProcessUserAttempt(ctx context.Context, event kafka.UserAttemptEvent) ([]kafka.UserProgressEvent, error) {
	var progressEvents []kafka.UserProgressEvent

	exercise, err := a.exerciseRepo.GetByQuestionID(ctx, event.QuestionUUID)
	if err != nil {
		return nil, err
	}

	if isCompleted, points := a.isExerciseCompleted(ctx, event.UserUUID, exercise.UUID, event.QuestionUUID, event.IsCorrect, event.SessionUUID); isCompleted {
		progressEvents = append(progressEvents, kafka.UserProgressEvent{
			UserUUID:   event.UserUUID,
			EntityType: "exercise",
			EntityUUID: exercise.UUID,
			Points:     points,
			IsCorrect:  true,
			CreatedAt:  time.Now(),
		})
	}

	if isCompleted, points := a.isLessonCompleted(ctx, event.UserUUID, exercise.LessonUUID, event.QuestionUUID, event.IsCorrect, event.SessionUUID); isCompleted {
		progressEvents = append(progressEvents, kafka.UserProgressEvent{
			UserUUID:   event.UserUUID,
			EntityType: "lesson",
			EntityUUID: exercise.LessonUUID,
			Points:     points,
			IsCorrect:  true,
			CreatedAt:  time.Now(),
		})
	}

	lesson, err := a.lessonRepo.GetByID(ctx, exercise.LessonUUID)
	if err != nil {
		return nil, err
	}

	if isCompleted, points := a.isCourseCompleted(ctx, event.UserUUID, lesson.CourseUUID, event.QuestionUUID, event.IsCorrect, event.SessionUUID); isCompleted {
		progressEvents = append(progressEvents, kafka.UserProgressEvent{
			UserUUID:   event.UserUUID,
			EntityType: "course",
			EntityUUID: lesson.CourseUUID,
			Points:     points,
			IsCorrect:  true,
			CreatedAt:  time.Now(),
		})
	}

	if len(progressEvents) == 0 {
		return nil, errors.New("no progress events generated")
	}

	return progressEvents, nil
}

func (a *AttemptService) isLessonCompleted(ctx context.Context, userUUID, lessonUUID, currentQuestionUUID uuid.UUID, currentIsCorrect bool, sessionUUID uuid.UUID) (bool, int) {
	completed, points, err := a.completionRepo.IsLessonCompleted(ctx, userUUID, lessonUUID)
	if err != nil || completed {
		return completed, points
	}

	exercises, err := a.exerciseRepo.GetByLessonID(ctx, lessonUUID)
	if err != nil {
		return false, 0
	}

	completedExercises := 0
	totalExercises := len(exercises)
	totalPoints := 0

	for _, exercise := range exercises {
		completed, exercisePoints, err := a.completionRepo.IsExerciseCompleted(ctx, userUUID, exercise.UUID)
		if err == nil && completed {
			completedExercises++
			totalPoints += exercisePoints
		}
	}

	if completedExercises == totalExercises {
		if err := a.completionRepo.MarkLessonCompleted(ctx, userUUID, lessonUUID, totalPoints/totalExercises); err != nil {
			return false, 0
		}
		return true, totalPoints / totalExercises
	}

	return false, 0
}

func (a *AttemptService) isCourseCompleted(ctx context.Context, userUUID, courseUUID, currentQuestionUUID uuid.UUID, currentIsCorrect bool, sessionUUID uuid.UUID) (bool, int) {
	completed, points, err := a.completionRepo.IsCourseCompleted(ctx, userUUID, courseUUID)
	if err != nil || completed {
		return completed, points
	}

	lessons, err := a.lessonRepo.GetByCourseID(ctx, courseUUID)
	if err != nil {
		return false, 0
	}

	completedLessons := 0
	totalLessons := len(lessons)
	totalPoints := 0

	for _, lesson := range lessons {
		completed, lessonPoints, err := a.completionRepo.IsLessonCompleted(ctx, userUUID, lesson.UUID)
		if err == nil && completed {
			completedLessons++
			totalPoints += lessonPoints
		}
	}

	if completedLessons == totalLessons {
		if err := a.completionRepo.MarkCourseCompleted(ctx, userUUID, courseUUID, totalPoints/totalLessons); err != nil {
			return false, 0
		}
		return true, totalPoints / totalLessons
	}

	return false, 0
}

func (a *AttemptService) isExerciseCompleted(ctx context.Context, userUUID, exerciseUUID, currentQuestionUUID uuid.UUID, currentIsCorrect bool, sessionUUID uuid.UUID) (bool, int) {
	totalQuestions, err := a.questionRepo.GetCountByExerciseID(ctx, exerciseUUID)
	if err != nil {
		return false, 0
	}

	exercise, err := a.exerciseRepo.GetByID(ctx, exerciseUUID)
	if err != nil {
		return false, 0
	}

	points := exercise.Points

	lastAttempts, err := a.attemptRepo.GetLast(ctx, userUUID, exerciseUUID, sessionUUID)
	if err != nil {
		return false, 0
	}

	correctQuestions := make(map[uuid.UUID]bool)

	for _, attempt := range lastAttempts {
		if attempt.IsCorrect {
			correctQuestions[attempt.QuestionUUID] = true
		}
	}

	if currentIsCorrect {
		correctQuestions[currentQuestionUUID] = true
	}

	if len(correctQuestions) < int(totalQuestions) {
		return false, 0
	}

	if totalQuestions == 0 {
		return false, 0
	}

	if err := a.completionRepo.MarkExerciseCompleted(ctx, userUUID, exerciseUUID, points); err != nil {
		return false, 0
	}

	return true, points

}
