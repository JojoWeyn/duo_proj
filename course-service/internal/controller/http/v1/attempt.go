package v1

import (
	"encoding/json"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type attemptRoutes struct {
	attemptUseCase  AttemptUseCase
	questionUseCase QuestionUseCase
}

func newAttemptRoutes(handler *gin.RouterGroup, attemptUseCase AttemptUseCase, questionUseCase QuestionUseCase) {
	r := &attemptRoutes{
		attemptUseCase:  attemptUseCase,
		questionUseCase: questionUseCase,
	}

	h := handler.Group("/attempts")
	{
		h.POST("/start/:exercise_id", r.startAttempt)
		h.POST("/answer", r.submitAnswer)
		h.POST("/finish", r.finishAttempt)
	}
}

func (a *attemptRoutes) startAttempt(c *gin.Context) {
	userUUID := c.GetHeader("X-User-UUID")
	exerciseUUID := c.Param("exercise_id")

	sessionID, err := a.attemptUseCase.StartAttempt(c.Request.Context(), uuid.MustParse(userUUID), uuid.MustParse(exerciseUUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.StartAttemptResponseDTO{SessionID: *sessionID})
}

func (a *attemptRoutes) submitAnswer(c *gin.Context) {
	sessionID := c.Query("session_id")
	userUUID := c.GetHeader("X-User-UUID")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing session_id in query"})
		return
	}

	var request struct {
		QuestionUUID string      `json:"question_id"`
		Answer       interface{} `json:"answer"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	correct, err := a.questionUseCase.CheckAnswer(c.Request.Context(), uuid.MustParse(userUUID), uuid.MustParse(request.QuestionUUID), request.Answer, uuid.MustParse(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := a.attemptUseCase.SubmitAnswer(c.Request.Context(), uuid.MustParse(sessionID), uuid.MustParse(userUUID), uuid.MustParse(request.QuestionUUID), convertAnswerToString(request.Answer), correct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SubmitAnswerResponseDTO{IsCorrect: correct})
}

func (a *attemptRoutes) finishAttempt(c *gin.Context) {
	sessionID := c.Query("session_id")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing session_id in query"})
		return
	}

	correctAnswers, questions, err := a.attemptUseCase.FinishAttempt(c.Request.Context(), uuid.MustParse(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var status bool
	if correctAnswers == questions {
		status = true
	} else {
		status = false
	}

	c.JSON(http.StatusOK, dto.FinishAttemptResponseDTO{
		IsFinished:     status,
		CorrectAnswers: correctAnswers,
		Questions:      questions,
	})
}

func convertAnswerToString(answer interface{}) string {
	jsonData, err := json.Marshal(answer)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}
