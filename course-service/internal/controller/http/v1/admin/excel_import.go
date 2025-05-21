package admin

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExcelImportUseCase интерфейс для импорта курсов из Excel
type ExcelImportUseCase interface {
	ImportCourseFromExcel(ctx context.Context, fileData []byte) error
}

// importCourseFromExcel обрабатывает запрос на импорт курса из Excel-файла
func (r *adminRoutes) importCourseFromExcel(c *gin.Context) {
	// Получаем файл из запроса
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл: " + err.Error()})
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения файла: " + err.Error()})
		return
	}

	// Импортируем курс из Excel
	if err := r.excelImportUseCase.ImportCourseFromExcel(c.Request.Context(), fileData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка импорта курса: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Курс успешно импортирован"})
}
