package main

import (
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/sqlite"
)

func main() {
	db, err := sqlite.NewSqliteDB(sqlite.Config{
		Path: "test.db",
	})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&entity.Course{}, &entity.Lesson{}, &entity.Question{}, &entity.QuestionOption{}, &entity.Exercise{}, &entity.Difficulty{}, &entity.CourseType{}); err != nil {
		panic(err)
	}

}
