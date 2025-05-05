package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type LessonUseCase interface {
	GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	GetLessonsByCourseID(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error)
}

type lessonRoutes struct {
	lessonUseCase LessonUseCase
}

func newLessonRoutes(handler *gin.RouterGroup, lessonUseCase LessonUseCase) {
	r := &lessonRoutes{
		lessonUseCase: lessonUseCase,
	}

	h := handler.Group("")
	{
		h.GET("/lesson/:id/info", r.getLessonByID)
		h.GET("/course/:id/content", r.getAllLessons)
	}
}

func (r *lessonRoutes) getLessonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	lesson, err := r.lessonUseCase.GetLessonByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var lessonFiles []dto.LessonFileDTO
	for _, file := range lesson.LessonFiles {
		lessonFiles = append(lessonFiles, dto.LessonFileDTO{
			UUID:       file.UUID,
			Title:      file.Title,
			FileURL:    file.FileURL,
			LessonUUID: file.LessonUUID,
		})
	}

	dto := dto.LessonDTO{
		UUID:         lesson.UUID,
		Title:        lesson.Title,
		Description:  lesson.Description,
		DifficultyID: lesson.DifficultyID,
		Order:        lesson.Order,
		CourseUUID:   lesson.CourseUUID,
		LessonFiles:  lessonFiles,
		Difficulty: dto.DifficultyDTO{
			ID:    lesson.Difficulty.ID,
			Title: lesson.Difficulty.Title,
		},
	}

	c.JSON(http.StatusOK, dto)
}

func (r *lessonRoutes) getAllLessons(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	lessons, err := r.lessonUseCase.GetLessonsByCourseID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var lessonDTOs []dto.LessonsDTO
	for _, lesson := range lessons {
		lessonDTOs = append(lessonDTOs, dto.LessonsDTO{
			UUID:         lesson.UUID,
			Title:        lesson.Title,
			Description:  lesson.Description,
			DifficultyID: lesson.DifficultyID,
			Order:        lesson.Order,
			CourseUUID:   lesson.CourseUUID,
			Difficulty: dto.DifficultyDTO{
				ID:    lesson.Difficulty.ID,
				Title: lesson.Difficulty.Title,
			},
		})
	}

	c.JSON(http.StatusOK, lessonDTOs)
}
