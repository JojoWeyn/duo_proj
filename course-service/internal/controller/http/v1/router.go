package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, au AttemptUseCase, gatewayUrl string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{gatewayUrl}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	handler.Use(cors.New(config))

	v1 := handler.Group("/v1")
	{
		newCourseRoutes(v1, cu)
		newLessonRoutes(v1, lu)
		newExerciseRoutes(v1, eu)
		newQuestionRoutes(v1, qu, au)
		newAttemptRoutes(v1, au, qu)
	}
}
