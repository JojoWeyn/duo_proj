package admin

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
)

type QuestionUseCase interface {
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*entity.Question, error)
	GetQuestionsByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]entity.Question, error)
	CreateQuestion(ctx context.Context, text string, typeID, order int, exerciseUUID uuid.UUID) error
	UpdateQuestion(ctx context.Context, question *entity.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
	AddImage(ctx context.Context, questionUUID uuid.UUID, title, fileUrl string) error
}

type ExerciseUseCase interface {
	GetExerciseByID(ctx context.Context, id uuid.UUID) (*entity.Exercise, error)
	GetExercisesByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*entity.Exercise, error)
	CreateExercise(ctx context.Context, title string, description string, points int, order int, lessonUUID uuid.UUID) error
	UpdateExercise(ctx context.Context, exercise *entity.Exercise) error
	DeleteExercise(ctx context.Context, id uuid.UUID) error
	AddFile(ctx context.Context, exerciseUUID uuid.UUID, title, fileUrl string) error
}

type LessonUseCase interface {
	GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	GetLessonsByCourseID(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, error)
	CreateLesson(ctx context.Context, title, description string, difficultyID, order int, courseUUID uuid.UUID) error
	UpdateLesson(ctx context.Context, lesson *entity.Lesson) error
	DeleteLesson(ctx context.Context, id uuid.UUID) error
}

type CourseUseCase interface {
	GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetAllCourses(ctx context.Context, typeId int) ([]entity.Course, error)
	CreateCourse(ctx context.Context, title, description string, typeID, difficultyID int) error
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error
}

type MatchingPairUseCase interface {
	GetMatchingPairByID(ctx context.Context, id uuid.UUID) (*entity.MatchingPair, error)
	GetMatchingPairsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.MatchingPair, error)
	CreateMatchingPair(ctx context.Context, left, right string, questionUUID uuid.UUID) error
	DeleteMatchingPair(ctx context.Context, id uuid.UUID) error
	UpdateMatchingPair(ctx context.Context, questionMatchingPair *entity.MatchingPair) error
}

type QuestionOptionUseCase interface {
	GetOptionByID(ctx context.Context, id uuid.UUID) (*entity.QuestionOption, error)
	GetOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]entity.QuestionOption, error)
	CreateQuestionOption(ctx context.Context, text string, isCorrect bool, questionUUID uuid.UUID) error
	DeleteQuestionOption(ctx context.Context, id uuid.UUID) error
	UpdateQuestionOption(ctx context.Context, option *entity.QuestionOption) error
}

type FileS3UseCase interface {
	UploadFile(ctx context.Context, file multipart.File, fileName string, fileSize int64, fileType string) (string, error)
	ListFiles(ctx context.Context) ([]string, error)
}

type adminRoutes struct {
	questionUseCase       QuestionUseCase
	exerciseUseCase       ExerciseUseCase
	lessonUseCase         LessonUseCase
	courseUseCase         CourseUseCase
	matchingPairUseCase   MatchingPairUseCase
	questionOptionUseCase QuestionOptionUseCase

	fileS3UseCase FileS3UseCase
}

func newAdminRoutes(handler *gin.RouterGroup, cu CourseUseCase, lu LessonUseCase, eu ExerciseUseCase, qu QuestionUseCase, mpu MatchingPairUseCase, qou QuestionOptionUseCase, fileS3UseCase FileS3UseCase) {
	r := &adminRoutes{
		questionUseCase:       qu,
		exerciseUseCase:       eu,
		lessonUseCase:         lu,
		courseUseCase:         cu,
		matchingPairUseCase:   mpu,
		questionOptionUseCase: qou,
		fileS3UseCase:         fileS3UseCase,
	}

	h := handler.Group("/admin")
	{
		h.GET("/course/list", r.getAllCourses)
		h.GET("/course/:course_id/lesson", r.getAllLessons)
		h.GET("/lesson/:lesson_id/exercise", r.getAllExercises)
		h.GET("/exercise/:exercise_id/question", r.getAllQuestions)
		h.GET("/question/:question_id/matching-pair", r.getAllMatchingPairs)
		h.GET("/question/:question_id/question-option", r.getAllQuestionOptions)

		h.GET("/course/:course_id/info", r.getCourseByID)
		h.GET("/lesson/:lesson_id/info", r.getLessonByID)
		h.GET("/exercise/:exercise_id/info", r.getExerciseByID)
		h.GET("/question/:question_id/info", r.getQuestionByID)
		h.GET("/matching-pair/:id/info", r.getMatchingPairByID)
		h.GET("/question-option/:id/info", r.getQuestionOptionByID)

		h.POST("/course", r.createCourse)
		h.POST("/lesson", r.createLesson)
		h.POST("/exercise", r.createExercise)
		h.POST("/question", r.createQuestion)
		h.POST("/matching-pair", r.createMatchingPair)
		h.POST("/question-option", r.createQuestionOption)

		h.POST("/file/add", r.addFile)

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

		h.POST("/file/upload", r.uploadFile)
		h.GET("/file/list", r.listFile)

	}
}

func (r *adminRoutes) addFile(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		FileUrl string `json:"file_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entityType := c.Query("entity")
	id := uuid.MustParse(c.Query("uuid"))

	switch entityType {
	case "exercise":
		if err := r.exerciseUseCase.AddFile(c.Request.Context(), id, req.Title, req.FileUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, nil)

	case "question":
		if err := r.questionUseCase.AddImage(c.Request.Context(), id, req.Title, req.FileUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, nil)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported entity type"})
	}

}

func (r *adminRoutes) listFile(c *gin.Context) {
	files, err := r.fileS3UseCase.ListFiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (r *adminRoutes) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileSize := file.Size
	fileName := file.Filename
	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileData.Close()

	buffer := make([]byte, 512)
	_, err = fileData.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file data"})
		return
	}

	fileType := http.DetectContentType(buffer)

	_, err = fileData.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset file reader"})
		return
	}

	fileURL, err := r.fileS3UseCase.UploadFile(c.Request.Context(), fileData, fileName, fileSize, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file on storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_url":  fileURL,
		"file_type": fileType,
	})
}

func (r *adminRoutes) getAllCourses(c *gin.Context) {
	courses, err := r.courseUseCase.GetAllCourses(c.Request.Context(), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (r *adminRoutes) getAllLessons(c *gin.Context) {
	id := c.Param("course_id")
	lessons, err := r.lessonUseCase.GetLessonsByCourseID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lessons)
}

func (r *adminRoutes) getAllExercises(c *gin.Context) {
	id := c.Param("lesson_id")
	exercises, err := r.exerciseUseCase.GetExercisesByLessonID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exercises)
}

func (r *adminRoutes) getAllQuestions(c *gin.Context) {
	id := c.Param("exercise_id")
	questions, err := r.questionUseCase.GetQuestionsByExerciseID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func (r *adminRoutes) getAllMatchingPairs(c *gin.Context) {
	id := c.Param("question_id")
	matchingPairs, err := r.matchingPairUseCase.GetMatchingPairsByQuestionID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matchingPairs)
}

func (r *adminRoutes) getAllQuestionOptions(c *gin.Context) {
	id := c.Param("question_id")
	questionOptions, err := r.questionOptionUseCase.GetOptionsByQuestionID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questionOptions)
}

func (r *adminRoutes) getCourseByID(c *gin.Context) {
	id := c.Param("course_id")
	course, err := r.courseUseCase.GetCourseByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (r *adminRoutes) getLessonByID(c *gin.Context) {
	id := c.Param("lesson_id")
	lesson, err := r.lessonUseCase.GetLessonByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lesson)
}

func (r *adminRoutes) getExerciseByID(c *gin.Context) {
	id := c.Param("exercise_id")
	exercise, err := r.exerciseUseCase.GetExerciseByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exercise)
}

func (r *adminRoutes) getQuestionByID(c *gin.Context) {
	id := c.Param("question_id")
	question, err := r.questionUseCase.GetQuestionByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, question)
}

func (r *adminRoutes) getMatchingPairByID(c *gin.Context) {
	id := c.Param("id")
	matchingPair, err := r.matchingPairUseCase.GetMatchingPairByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matchingPair)
}

func (r *adminRoutes) getQuestionOptionByID(c *gin.Context) {
	id := c.Param("id")
	questionOption, err := r.questionOptionUseCase.GetOptionByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questionOption)
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
