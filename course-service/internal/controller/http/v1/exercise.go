package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
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
	lessonID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	exercises, err := h.exerciseUseCase.GetExercisesByLessonID(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var dtoList []dto.ExerciseAllDTO
	for _, ex := range exercises {
		dtoList = append(dtoList, dto.ExerciseAllDTO{
			UUID:        ex.UUID,
			Title:       ex.Title,
			Description: ex.Description,
			Points:      ex.Points,
			Order:       ex.Order,
			LessonUUID:  ex.LessonUUID,
		})
	}

	c.JSON(http.StatusOK, dtoList)
}

func (h *exerciseRoutes) getExerciseByID(c *gin.Context) {
	id := c.Param("id")
	exerciseID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	exercise, err := h.exerciseUseCase.GetExerciseByID(c.Request.Context(), exerciseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var fileDTOs []dto.ExerciseFileDTO
	for _, file := range exercise.ExerciseFiles {
		fileDTOs = append(fileDTOs, dto.ExerciseFileDTO{
			UUID:         file.UUID,
			Title:        file.Title,
			FileURL:      file.FileURL,
			ExerciseUUID: file.ExerciseUUID,
		})
	}

	// Формируем DTO для ответа
	exerciseDTO := dto.ExerciseOneDTO{
		UUID:          exercise.UUID,
		Title:         exercise.Title,
		Description:   exercise.Description,
		Points:        exercise.Points,
		Order:         exercise.Order,
		LessonUUID:    exercise.LessonUUID,
		ExerciseFiles: fileDTOs,
	}

	c.JSON(http.StatusOK, exerciseDTO)
}
