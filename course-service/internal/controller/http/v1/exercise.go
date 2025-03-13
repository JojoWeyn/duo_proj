package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ExerciseUseCase interface {
	GetExercisesByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error)
	GetExerciseByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error)
	UpdateExercise(ctx context.Context, exercise *entity.Exercise) error
	DeleteExercise(ctx context.Context, id uuid.UUID) error
}
type exerciseRoutes struct {
	exerciseUseCase ExerciseUseCase
}

func newExerciseRoutes(handler *gin.RouterGroup, exerciseUseCase ExerciseUseCase) {
	r := &exerciseRoutes{
		exerciseUseCase: exerciseUseCase,
	}

	h := handler.Group("")
	{
		h.GET("/lesson/:id/content", r.getAllExercises)
		h.GET("/exercise/:id/info", r.getExerciseByID)
	}
}

func (h *exerciseRoutes) getAllExercises(c *gin.Context) {
	id := c.Param("id")

	exercises, err := h.exerciseUseCase.GetExercisesByLessonID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, exercises)
}

func (h *exerciseRoutes) getExerciseByID(c *gin.Context) {
	id := c.Param("id")
	exercise, err := h.exerciseUseCase.GetExerciseByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, exercise)
}
