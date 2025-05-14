package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type CourseUseCase interface {
	GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAllCourses(ctx context.Context, title string, diffId int, typeId int) ([]entity.Course, error)
}

type courseRoutes struct {
	courseUseCase CourseUseCase
}

func newCourseRoutes(handler *gin.RouterGroup, courseUseCase CourseUseCase) {
	r := &courseRoutes{
		courseUseCase: courseUseCase,
	}

	h := handler.Group("/course")
	{
		h.GET("/:id/info", r.getCourseByID)
		h.GET("/list", r.getAllCourses)
	}
}

func (r *courseRoutes) getCourseByID(c *gin.Context) {
	id := c.Param("id")
	course, err := r.courseUseCase.GetCourseByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	courseDTO := dto.CourseInfoDTO{
		UUID:         course.UUID,
		Title:        course.Title,
		Description:  course.Description,
		DifficultyID: course.DifficultyID,
		TypeID:       course.TypeID,
		CourseType:   course.CourseType,
		Difficulty:   course.Difficulty,
	}

	c.JSON(http.StatusOK, courseDTO)
}

func (r *courseRoutes) getAllCourses(c *gin.Context) {
	title := c.Query("title")
	diffIdStr := c.Query("difficultyId")
	diffId, _ := strconv.Atoi(diffIdStr)

	courses, err := r.courseUseCase.GetAllCourses(c.Request.Context(), title, diffId, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coursesDTO := make([]dto.CourseSmallDTO, 0, len(courses))
	for _, course := range courses {
		coursesDTO = append(coursesDTO, dto.CourseSmallDTO{
			UUID:         course.UUID,
			Title:        course.Title,
			Description:  course.Description,
			DifficultyID: course.DifficultyID,
			Difficulty: dto.DifficultyDTO{
				ID:    course.Difficulty.ID,
				Title: course.Difficulty.Title,
			},
		})
	}

	c.JSON(http.StatusOK, coursesDTO)
}
