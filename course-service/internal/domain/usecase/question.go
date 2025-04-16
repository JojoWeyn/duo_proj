package usecase

import (
	"context"
	"errors"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"log"
)

type QuestionRepository interface {
	AddImage(ctx context.Context, file *entity.QuestionImage) error
	Create(ctx context.Context, question *entity.Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)
	GetByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error)
	Update(ctx context.Context, question *entity.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	ProcessUserAttempt(ctx context.Context, event kafka.UserAttemptEvent) ([]kafka.UserProgressEvent, error)
}

type QuestionUseCase struct {
	repo     QuestionRepository
	producer *kafka.Producer
	service  Service
}

func NewQuestionUseCase(repo QuestionRepository, producer *kafka.Producer, service Service) *QuestionUseCase {
	return &QuestionUseCase{
		repo:     repo,
		producer: producer,
		service:  service,
	}
}

func (q *QuestionUseCase) CheckAnswer(ctx context.Context, userUUID, questionID uuid.UUID, userAnswers interface{}, sessionUUID uuid.UUID) (bool, error) {
	question, err := q.repo.GetByID(ctx, questionID)
	if err != nil {
		return false, err
	}

	isCorrect, err := func() (bool, error) {
		switch question.TypeID {
		case 1: // Single Choice
			userAnswerUUID, ok := userAnswers.(string)
			if !ok {
				return false, errors.New("1invalid answer format")
			}
			for _, option := range question.QuestionOptions {
				if option.IsCorrect && option.UUID == uuid.MustParse(userAnswerUUID) {
					return true, nil
				}
			}
			return false, nil

		case 2: // Multiple Choice
			userAnswerUUIDs, ok := userAnswers.([]interface{})
			if !ok {
				return false, errors.New("2invalid answer format")
			}

			var answers []uuid.UUID
			for _, a := range userAnswerUUIDs {
				answerStr, ok := a.(string)
				if !ok {
					return false, errors.New("2invalid answer format")
				}
				answers = append(answers, uuid.MustParse(answerStr))
			}

			correctAnswers := make(map[uuid.UUID]bool)
			for _, option := range question.QuestionOptions {
				if option.IsCorrect {
					correctAnswers[option.UUID] = true
				}
			}
			if len(correctAnswers) != len(answers) {
				return false, nil
			}

			for _, answer := range answers {
				if !correctAnswers[answer] {
					return false, nil
				}
			}
			return true, nil

		case 3: // Matching
			userPairs, ok := userAnswers.(map[string]interface{}) // map[left] = right
			if !ok {
				return false, errors.New("3 invalid answer format")
			}

			userAnswerMap := make(map[string]string)
			for left, right := range userPairs {
				rightStr, ok := right.(string)
				if !ok {
					return false, errors.New("3 invalid answer format")
				}
				userAnswerMap[left] = rightStr
			}

			correctPairs := make(map[string]string)
			for _, pair := range question.MatchingPairs {
				correctPairs[pair.LeftText] = pair.RightText
			}
			if len(correctPairs) != len(userAnswerMap) {
				return false, nil
			}
			for left, right := range userAnswerMap {
				if correctPairs[left] != right {
					return false, nil
				}
			}
			return true, nil
		}

		return false, errors.New("unknown question type")
	}()

	if err == nil {
		go func() {
			if err := q.producer.SendUserAttemptEvent(userUUID, questionID, isCorrect, sessionUUID); err != nil {
				log.Printf("Failed to send progress event: %v", err)
			}

			if isCorrect == true {
				if err := q.producer.SendUserEvent(kafka.UserEvent{
					UUID:   userUUID.String(),
					Login:  "",
					Action: "question",
				}); err != nil {
					log.Printf("Failed to send progress event: %v", err)
				}
			}

		}()

		userProgressEvents, err := q.service.ProcessUserAttempt(ctx, kafka.UserAttemptEvent{
			UserUUID:     userUUID,
			QuestionUUID: questionID,
			SessionUUID:  sessionUUID,
			IsCorrect:    isCorrect,
		})
		if err != nil {
			log.Printf("Failed to process user attempt: %v", err)
		}
		if len(userProgressEvents) > 0 {
			go func() {
				for _, event := range userProgressEvents {
					if err := q.producer.SendUserProgressEvent(event); err != nil {
						log.Printf("Failed to send progress event: %v", err)
					}
					if err := q.producer.SendUserEvent(kafka.UserEvent{
						UUID:   event.UserUUID.String(),
						Login:  "",
						Action: event.EntityType,
					}); err != nil {
						log.Printf("Failed to send progress event: %v", err)
					}

				}
			}()
		}
	}

	return isCorrect, err

}

func (q *QuestionUseCase) CreateQuestion(ctx context.Context, text string, typeID, order int, exerciseUUID uuid.UUID) error {
	question := entity.NewQuestion(text, typeID, order, exerciseUUID)
	return q.repo.Create(ctx, question)
}

func (q *QuestionUseCase) GetQuestionByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	return q.repo.GetByID(ctx, id)
}

func (q *QuestionUseCase) GetQuestionsByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error) {
	return q.repo.GetByExerciseID(ctx, exerciseID)
}

func (q *QuestionUseCase) UpdateQuestion(ctx context.Context, question *entity.Question) error {
	return q.repo.Update(ctx, question)
}

func (q *QuestionUseCase) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return q.repo.Delete(ctx, id)
}

func (q *QuestionUseCase) AddImage(ctx context.Context, questionUUID uuid.UUID, title, fileUrl string) error {
	file := entity.QuestionImage{
		UUID:         uuid.New(),
		Title:        title,
		ImageURL:     fileUrl,
		QuestionUUID: questionUUID,
	}

	return q.repo.AddImage(ctx, &file)
}

func (q *QuestionUseCase) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return q.repo.DeleteImage(ctx, id)
}
