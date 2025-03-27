package admin

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, mpu MatchingPairUseCase, qou QuestionOptionUseCase, gatewayUrl string) {

	v1 := handler.Group("/v1")
	{
		newAdminRoutes(v1, cu, lu, eu, qu, mpu, qou)
	}
}
