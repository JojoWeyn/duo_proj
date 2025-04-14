package admin

import (
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, mpu MatchingPairUseCase, qou QuestionOptionUseCase, fileS3UseCase FileS3UseCase) {

	v1 := handler.Group("/v1", middleware.RoleMiddleware("admin"))
	{
		newAdminRoutes(v1, cu, lu, eu, qu, mpu, qou, fileS3UseCase)
	}
}
