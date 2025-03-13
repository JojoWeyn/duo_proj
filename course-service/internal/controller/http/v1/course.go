package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CourseUseCase interface {
	CreateCourse(ctx context.Context, title, description string, typeID int) error
	GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAllCourses(ctx context.Context) ([]*entity.Course, error)
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error
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
	courses, err := r.courseUseCase.GetAllCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var coursesDTO []dto.CourseSmallDTO
	for _, course := range courses {
		coursesDTO = append(coursesDTO, dto.CourseSmallDTO{
			UUID:         course.UUID,
			Title:        course.Title,
			Description:  course.Description,
			DifficultyID: course.DifficultyID,
			Difficulty:   course.Difficulty,
		})
	}

	c.JSON(http.StatusOK, coursesDTO)
}
