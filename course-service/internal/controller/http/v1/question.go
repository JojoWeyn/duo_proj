package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"time"
)

type QuestionUseCase interface {
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)
	GetQuestionsByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error)
	CheckAnswer(ctx context.Context, userUUID, questionID uuid.UUID, userAnswers interface{}, sessionUUID uuid.UUID) (bool, error)
}

type AttemptUseCase interface {
	StartAttempt(ctx context.Context, userUUID, exerciseUUID uuid.UUID) (*uuid.UUID, error)
	SubmitAnswer(ctx context.Context, sessionID, userUUID, questionUUID uuid.UUID, answer string, isCorrect bool) error
	FinishAttempt(ctx context.Context, sessionID uuid.UUID) (int, int, error)
	CreateAttempt(ctx context.Context, userUUID, questionUUID uuid.UUID, answer string, isCorrect bool) error
	GetAttemptsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Attempt, error)
}

type questionRoutes struct {
	questionUseCase QuestionUseCase
	attemptUseCase  AttemptUseCase
}

func newQuestionRoutes(handler *gin.RouterGroup, questionUseCase QuestionUseCase, attemptUseCase AttemptUseCase) {
	r := &questionRoutes{
		questionUseCase: questionUseCase,
		attemptUseCase:  attemptUseCase,
	}

	h := handler.Group("")
	{
		h.GET("/question/:id/info", r.getQuestionByID)
		h.GET("/exercise/:id/question", r.getAllQuestions)
	}
}

func (q *questionRoutes) getQuestionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	question, err := q.questionUseCase.GetQuestionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var imagesDTO []dto.QuestionImageDTO
	for _, image := range question.QuestionImages {
		imagesDTO = append(imagesDTO, dto.QuestionImageDTO{
			ImageURL: image.ImageURL,
		})
	}

	var optionsDTO []dto.QuestionOptionDTO
	for _, option := range question.QuestionOptions {
		optionsDTO = append(optionsDTO, dto.QuestionOptionDTO{
			UUID: option.UUID,
			Text: option.Text,
		})
	}

	var leftSide, rightSide []string
	for _, pair := range question.MatchingPairs {
		leftSide = append(leftSide, pair.LeftText)
		rightSide = append(rightSide, pair.RightText)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rightSide), func(i, j int) {
		rightSide[i], rightSide[j] = rightSide[j], rightSide[i]
	})

	var matchingDTO dto.QuestionMatchingDTO
	matchingDTO.LeftSide = leftSide
	matchingDTO.RightSide = rightSide

	questionDTO := dto.QuestionDTO{
		UUID:            question.UUID,
		TypeID:          question.TypeID,
		QuestionType:    question.QuestionType,
		Text:            question.Text,
		Images:          imagesDTO,
		Order:           question.Order,
		ExerciseUUID:    question.ExerciseUUID,
		QuestionOptions: optionsDTO,
		Matching:        matchingDTO,
	}

	c.JSON(http.StatusOK, questionDTO)
}

func (q *questionRoutes) getAllQuestions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	questions, err := q.questionUseCase.GetQuestionsByExerciseID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var questionDTOs []dto.QuestionsDTO
	for _, question := range questions {
		questionDTOs = append(questionDTOs, dto.QuestionsDTO{
			UUID:         question.UUID,
			TypeID:       question.TypeID,
			Text:         question.Text,
			Order:        question.Order,
			ExerciseUUID: question.ExerciseUUID,
			Type: dto.QuestionTypeDTO{
				ID:    question.QuestionType.ID,
				Title: question.QuestionType.Title,
			},
		})
	}

	c.JSON(http.StatusOK, questionDTOs)
}
