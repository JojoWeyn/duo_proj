package admin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, mpu MatchingPairUseCase, qou QuestionOptionUseCase, gatewayUrl string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	handler.Use(cors.New(config))

	v1 := handler.Group("/v1")
	{
		newAdminRoutes(v1, cu, lu, eu, qu, mpu, qou)
	}
}
