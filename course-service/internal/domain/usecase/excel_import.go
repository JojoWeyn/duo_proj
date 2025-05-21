package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// ExcelImportUseCase представляет сервис для импорта курсов из Excel-файлов
type ExcelImportUseCase struct {
	courseRepo         CourseRepository
	lessonRepo         LessonRepository
	exerciseRepo       ExerciseRepository
	questionRepo       QuestionRepository
	matchingPairRepo   MatchingPairRepository
	questionOptionRepo QuestionOptionRepository
}

// NewExcelImportUseCase создает новый экземпляр сервиса импорта из Excel
func NewExcelImportUseCase(
	courseRepo CourseRepository,
	lessonRepo LessonRepository,
	exerciseRepo ExerciseRepository,
	questionRepo QuestionRepository,
	matchingPairRepo MatchingPairRepository,
	questionOptionRepo QuestionOptionRepository,
) *ExcelImportUseCase {
	return &ExcelImportUseCase{
		courseRepo:         courseRepo,
		lessonRepo:         lessonRepo,
		exerciseRepo:       exerciseRepo,
		questionRepo:       questionRepo,
		matchingPairRepo:   matchingPairRepo,
		questionOptionRepo: questionOptionRepo,
	}
}

// ImportCourseFromExcel импортирует курс из Excel-файла
func (e *ExcelImportUseCase) ImportCourseFromExcel(ctx context.Context, fileData []byte) error {
	// Открываем Excel-файл из байтов
	f, err := excelize.OpenReader(strings.NewReader(string(fileData)))
	if err != nil {
		return fmt.Errorf("ошибка открытия Excel-файла: %w", err)
	}
	defer f.Close()

	// Получаем список листов
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return errors.New("Excel-файл не содержит листов")
	}

	// Обрабатываем первый лист как информацию о курсе
	courseSheet := sheets[0]
	course, err := e.parseCourseSheet(f, courseSheet)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге информации о курсе: %w", err)
	}

	// Создаем курс в базе данных
	if err := e.courseRepo.Create(ctx, course); err != nil {
		return fmt.Errorf("ошибка при создании курса: %w", err)
	}

	// Обрабатываем остальные листы как уроки
	for i := 1; i < len(sheets); i++ {
		lessonSheet := sheets[i]
		lesson, err := e.parseLessonSheet(f, lessonSheet, course.UUID)
		if err != nil {
			return fmt.Errorf("ошибка при парсинге урока %s: %w", lessonSheet, err)
		}

		// Создаем урок в базе данных
		if err := e.lessonRepo.Create(ctx, lesson); err != nil {
			return fmt.Errorf("ошибка при создании урока: %w", err)
		}

		// Обрабатываем упражнения для урока
		if err := e.parseExercises(ctx, f, lessonSheet, lesson.UUID); err != nil {
			return fmt.Errorf("ошибка при парсинге упражнений для урока %s: %w", lessonSheet, err)
		}
	}

	return nil
}

// parseCourseSheet парсит информацию о курсе из листа Excel
func (e *ExcelImportUseCase) parseCourseSheet(f *excelize.File, sheetName string) (*entity.Course, error) {
	// Получаем заголовок курса (ячейка A1)
	title, err := f.GetCellValue(sheetName, "A2")
	if err != nil || title == "" {
		return nil, errors.New("не удалось получить заголовок курса")
	}

	// Получаем описание курса (ячейка A2)
	description, _ := f.GetCellValue(sheetName, "B2")

	// Получаем тип курса (ячейка B1)
	typeIDStr, _ := f.GetCellValue(sheetName, "C2")
	typeID, err := strconv.Atoi(typeIDStr)
	if err != nil {
		typeID = 1 // Значение по умолчанию, если не указано
	}

	// Получаем сложность курса (ячейка B2)
	difficultyIDStr, _ := f.GetCellValue(sheetName, "D2")
	difficultyID, err := strconv.Atoi(difficultyIDStr)
	if err != nil {
		difficultyID = 1 // Значение по умолчанию, если не указано
	}

	// Создаем новый курс
	return entity.NewCourse(title, description, typeID, difficultyID), nil
}

// parseLessonSheet парсит информацию об уроке из листа Excel
func (e *ExcelImportUseCase) parseLessonSheet(f *excelize.File, sheetName string, courseUUID uuid.UUID) (*entity.Lesson, error) {
	// Получаем заголовок урока (ячейка A1)
	title, err := f.GetCellValue(sheetName, "A2")
	if err != nil || title == "" {
		return nil, errors.New("не удалось получить заголовок урока")
	}

	// Получаем описание урока (ячейка A2)
	description, _ := f.GetCellValue(sheetName, "B2")

	// Получаем сложность урока (ячейка B1)
	difficultyIDStr, _ := f.GetCellValue(sheetName, "C2")
	difficultyID, err := strconv.Atoi(difficultyIDStr)
	if err != nil {
		difficultyID = 1 // Значение по умолчанию, если не указано
	}

	// Получаем порядковый номер урока (ячейка B2)
	orderStr, _ := f.GetCellValue(sheetName, "D2")
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		order = 1 // Значение по умолчанию, если не указано
	}

	// Создаем новый урок
	return entity.NewLesson(title, description, difficultyID, order, courseUUID), nil
}

// parseExercises парсит упражнения из листа Excel
func (e *ExcelImportUseCase) parseExercises(ctx context.Context, f *excelize.File, sheetName string, lessonUUID uuid.UUID) error {
	// Начинаем с 6-й строки (после заголовка и описания урока)
	row := 6
	exerciseOrder := 1

	for {
		// Проверяем, есть ли заголовок упражнения
		exerciseCell := fmt.Sprintf("A%d", row)
		exerciseTitle, err := f.GetCellValue(sheetName, exerciseCell)
		if err != nil || exerciseTitle == "" {
			// Больше нет упражнений
			break
		}

		// Получаем описание упражнения
		descriptionCell := fmt.Sprintf("B%d", row)
		exerciseDescription, _ := f.GetCellValue(sheetName, descriptionCell)

		// Получаем количество баллов за упражнение
		pointsCell := fmt.Sprintf("C%d", row)
		pointsStr, _ := f.GetCellValue(sheetName, pointsCell)
		points, err := strconv.Atoi(pointsStr)
		if err != nil {
			points = 10 // Значение по умолчанию, если не указано
		}

		// Создаем новое упражнение
		exercise := entity.NewExercise(exerciseTitle, exerciseDescription, points, exerciseOrder, lessonUUID)
		if err := e.exerciseRepo.Create(ctx, exercise); err != nil {
			return fmt.Errorf("ошибка при создании упражнения: %w", err)
		}

		// Переходим к вопросам (начиная с row+4)
		questionRow := row + 4
		questionOrder := 1

		for {
			// Проверяем, есть ли текст вопроса
			questionCell := fmt.Sprintf("A%d", questionRow)
			questionText, err := f.GetCellValue(sheetName, questionCell)
			if err != nil || questionText == "" {
				// Больше нет вопросов или начинается новое упражнение
				break
			}

			// Получаем тип вопроса
			typeCell := fmt.Sprintf("B%d", questionRow)
			typeIDStr, _ := f.GetCellValue(sheetName, typeCell)
			typeID, err := strconv.Atoi(typeIDStr)
			if err != nil {
				typeID = 1 // Значение по умолчанию (одиночный выбор)
			}

			// Создаем новый вопрос
			question := entity.NewQuestion(questionText, typeID, questionOrder, exercise.UUID)
			if err := e.questionRepo.Create(ctx, question); err != nil {
				return fmt.Errorf("ошибка при создании вопроса: %w", err)
			}

			// Обрабатываем варианты ответов или пары для сопоставления
			if typeID == 3 { // Тип вопроса - сопоставление
				if err := e.parseMatchingPairs(ctx, f, sheetName, questionRow+2, question.UUID); err != nil {
					return err
				}
			} else { // Одиночный или множественный выбор
				if err := e.parseQuestionOptions(ctx, f, sheetName, questionRow+2, question.UUID); err != nil {
					return err
				}
			}

			// Переходим к следующему вопросу
			questionRow += 7 // Примерно 5 строк на вопрос и варианты ответов
			questionOrder++
		}

		// Переходим к следующему упражнению
		row = questionRow + 2 // Пропускаем пустую строку между упражнениями
		exerciseOrder++
	}

	return nil
}

// parseQuestionOptions парсит варианты ответов для вопроса
func (e *ExcelImportUseCase) parseQuestionOptions(ctx context.Context, f *excelize.File, sheetName string, startRow int, questionUUID uuid.UUID) error {
	for i := 0; i < 4; i++ { // Максимум 4 варианта ответа
		row := startRow + i

		// Получаем текст варианта ответа
		optionCell := fmt.Sprintf("A%d", row)
		optionText, err := f.GetCellValue(sheetName, optionCell)
		if err != nil || optionText == "" {
			continue // Пропускаем пустые варианты
		}

		// Получаем признак правильного ответа
		correctCell := fmt.Sprintf("B%d", row)
		correctStr, _ := f.GetCellValue(sheetName, correctCell)
		isCorrect := strings.ToLower(correctStr) == "true" || correctStr == "1"

		// Создаем новый вариант ответа
		option := entity.NewQuestionOption(optionText, isCorrect, questionUUID)
		if err := e.questionOptionRepo.Create(ctx, option); err != nil {
			return fmt.Errorf("ошибка при создании варианта ответа: %w", err)
		}
	}

	return nil
}

// parseMatchingPairs парсит пары для сопоставления
func (e *ExcelImportUseCase) parseMatchingPairs(ctx context.Context, f *excelize.File, sheetName string, startRow int, questionUUID uuid.UUID) error {
	for i := 0; i < 4; i++ { // Максимум 4 пары для сопоставления
		row := startRow + i

		// Получаем левую часть пары
		leftCell := fmt.Sprintf("A%d", row)
		leftText, err := f.GetCellValue(sheetName, leftCell)
		if err != nil || leftText == "" {
			continue // Пропускаем пустые пары
		}

		// Получаем правую часть пары
		rightCell := fmt.Sprintf("B%d", row)
		rightText, _ := f.GetCellValue(sheetName, rightCell)

		// Создаем новую пару для сопоставления
		pair := entity.NewMatchingPair(leftText, rightText, questionUUID)
		if err := e.matchingPairRepo.Create(ctx, pair); err != nil {
			return fmt.Errorf("ошибка при создании пары для сопоставления: %w", err)
		}
	}

	return nil
}
