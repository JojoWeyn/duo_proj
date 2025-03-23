package v1

import (
	"encoding/json"
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
		h.POST("/:session_id/answer", r.submitAnswer)
		h.POST("/:session_id/finish", r.finishAttempt)
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

	c.JSON(http.StatusOK, gin.H{"session_id": sessionID})
}

func (a *attemptRoutes) submitAnswer(c *gin.Context) {
	sessionID := c.Param("session_id")
	userUUID := c.GetHeader("X-User-UUID")

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

	c.JSON(http.StatusOK, gin.H{"is_correct": correct})
}

func (a *attemptRoutes) finishAttempt(c *gin.Context) {
	sessionID := c.Param("session_id")

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

	c.JSON(http.StatusOK, gin.H{
		"is_finished":     status,
		"correct_answers": correctAnswers,
		"questions":       questions,
	})
}

func convertAnswerToString(answer interface{}) string {
	jsonData, err := json.Marshal(answer)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}
