package admin

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type QuestionUseCase interface {
	CreateQuestion(ctx context.Context, text string, typeID, order int, exerciseUUID uuid.UUID) error
	UpdateQuestion(ctx context.Context, question *entity.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
}

type ExerciseUseCase interface {
	CreateExercise(ctx context.Context, title string, description string, points int, order int, lessonUUID uuid.UUID) error
	UpdateExercise(ctx context.Context, exercise *entity.Exercise) error
	DeleteExercise(ctx context.Context, id uuid.UUID) error
}

type LessonUseCase interface {
	CreateLesson(ctx context.Context, title, description string, difficultyID, order int, courseUUID uuid.UUID) error
	UpdateLesson(ctx context.Context, lesson *entity.Lesson) error
	DeleteLesson(ctx context.Context, id uuid.UUID) error
}

type CourseUseCase interface {
	CreateCourse(ctx context.Context, title, description string, typeID, difficultyID int) error
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error
}

type MatchingPairUseCase interface {
	CreateMatchingPair(ctx context.Context, left, right string, questionUUID uuid.UUID) error
	DeleteMatchingPair(ctx context.Context, id uuid.UUID) error
	UpdateMatchingPair(ctx context.Context, questionMatchingPair *entity.MatchingPair) error
}

type QuestionOptionUseCase interface {
	CreateQuestionOption(ctx context.Context, text string, isCorrect bool, questionUUID uuid.UUID) error
	DeleteQuestionOption(ctx context.Context, id uuid.UUID) error
	UpdateQuestionOption(ctx context.Context, option *entity.QuestionOption) error
}

type adminRoutes struct {
	questionUseCase       QuestionUseCase
	exerciseUseCase       ExerciseUseCase
	lessonUseCase         LessonUseCase
	courseUseCase         CourseUseCase
	matchingPairUseCase   MatchingPairUseCase
	questionOptionUseCase QuestionOptionUseCase
}

func newAdminRoutes(handler *gin.RouterGroup, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, mpu MatchingPairUseCase, qou QuestionOptionUseCase) {
	r := &adminRoutes{
		questionUseCase:       qu,
		exerciseUseCase:       eu,
		lessonUseCase:         lu,
		courseUseCase:         cu,
		matchingPairUseCase:   mpu,
		questionOptionUseCase: qou,
	}

	h := handler.Group("/admin")
	{
		h.POST("/course", r.createCourse)
		h.POST("/lesson", r.createLesson)
		h.POST("/exercise", r.createExercise)
		h.POST("/question", r.createQuestion)
		h.POST("/matching-pair", r.createMatchingPair)
		h.POST("/question-option", r.createQuestionOption)

		h.PATCH("/course/:id", r.updateCourse)
		h.PATCH("/lesson/:id", r.updateLesson)
		h.PATCH("/exercise/:id", r.updateExercise)
		h.PATCH("/question/:id", r.updateQuestion)
		h.PATCH("/matching-pair/:id", r.updateMatchingPair)
		h.PATCH("/question-option/:id", r.updateQuestionOption)

		h.DELETE("/course/:id", r.deleteCourse)
		h.DELETE("/lesson/:id", r.deleteLesson)
		h.DELETE("/exercise/:id", r.deleteExercise)
		h.DELETE("/question/:id", r.deleteQuestion)
		h.DELETE("/matching-pair/:id", r.deleteMatchingPair)
		h.DELETE("/question-option/:id", r.deleteQuestionOption)

	}
}

func (r *adminRoutes) createCourse(c *gin.Context) {
	var course dto.CourseCreateRequestDTO
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.courseUseCase.CreateCourse(c.Request.Context(), course.Title, course.Description, course.TypeID, course.DifficultyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (r *adminRoutes) createLesson(c *gin.Context) {
	var lesson dto.LessonCreateRequestDTO
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.lessonUseCase.CreateLesson(c.Request.Context(), lesson.Title, lesson.Description, lesson.DifficultyID, lesson.Order, lesson.CourseUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

func (r *adminRoutes) createExercise(c *gin.Context) {
	var exercise dto.ExerciseCreateRequestDTO
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.exerciseUseCase.CreateExercise(c.Request.Context(), exercise.Title, exercise.Description, exercise.Points, exercise.Order, exercise.LessonUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (r *adminRoutes) createQuestion(c *gin.Context) {
	var question dto.QuestionCreateRequestDTO
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.questionUseCase.CreateQuestion(c.Request.Context(), question.Text, question.TypeID, question.Order, question.ExerciseUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

func (r *adminRoutes) createMatchingPair(c *gin.Context) {
	var matchingPair dto.MatchingPairCreateRequestDTO
	if err := c.ShouldBindJSON(&matchingPair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.matchingPairUseCase.CreateMatchingPair(c.Request.Context(), matchingPair.LeftText, matchingPair.RightText, matchingPair.QuestionUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, matchingPair)
}

func (r *adminRoutes) createQuestionOption(c *gin.Context) {
	var questionOption dto.QuestionOptionCreateRequestDTO
	if err := c.ShouldBindJSON(&questionOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.questionOptionUseCase.CreateQuestionOption(c.Request.Context(), questionOption.Text, questionOption.IsCorrect, questionOption.QuestionUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, questionOption)
}

func (r *adminRoutes) updateCourse(c *gin.Context) {
	var course entity.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.courseUseCase.UpdateCourse(c.Request.Context(), &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (r *adminRoutes) updateLesson(c *gin.Context) {
	var lesson entity.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.lessonUseCase.UpdateLesson(c.Request.Context(), &lesson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (r *adminRoutes) updateExercise(c *gin.Context) {
	var exercise entity.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.exerciseUseCase.UpdateExercise(c.Request.Context(), &exercise); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (r *adminRoutes) updateQuestion(c *gin.Context) {
	var question entity.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.questionUseCase.UpdateQuestion(c.Request.Context(), &question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (r *adminRoutes) updateMatchingPair(c *gin.Context) {
	var matchingPair entity.MatchingPair
	if err := c.ShouldBindJSON(&matchingPair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.matchingPairUseCase.UpdateMatchingPair(c.Request.Context(), &matchingPair); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matchingPair)
}

func (r *adminRoutes) updateQuestionOption(c *gin.Context) {
	var questionOption entity.QuestionOption
	if err := c.ShouldBindJSON(&questionOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.questionOptionUseCase.UpdateQuestionOption(c.Request.Context(), &questionOption); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questionOption)
}

func (r *adminRoutes) deleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := r.courseUseCase.DeleteCourse(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *adminRoutes) deleteLesson(c *gin.Context) {
	id := c.Param("id")
	if err := r.lessonUseCase.DeleteLesson(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *adminRoutes) deleteExercise(c *gin.Context) {
	id := c.Param("id")
	if err := r.exerciseUseCase.DeleteExercise(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *adminRoutes) deleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if err := r.questionUseCase.DeleteQuestion(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *adminRoutes) deleteMatchingPair(c *gin.Context) {
	id := c.Param("id")
	if err := r.matchingPairUseCase.DeleteMatchingPair(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *adminRoutes) deleteQuestionOption(c *gin.Context) {
	id := c.Param("id")
	if err := r.questionOptionUseCase.DeleteQuestionOption(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
