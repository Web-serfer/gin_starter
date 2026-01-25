package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"gin-starter/internal/service/image"

	"github.com/gin-gonic/gin"
)

// ImageHandler структура для обработки запросов к изображениям
type ImageHandler struct {
	processor *image.ProcessorService
}

// NewImageHandler создает новый экземпляр ImageHandler
func NewImageHandler(processor *image.ProcessorService) *ImageHandler {
	return &ImageHandler{
		processor: processor,
	}
}

// OptimizedImage обрабатывает запросы на оптимизацию изображений
func (ih *ImageHandler) OptimizedImage(c *gin.Context) {
	path := c.Query("path")
	widthStr := c.Query("w")
	heightStr := c.Query("h")
	qualityStr := c.Query("q")

	// Проверяем обязательный параметр path
	if path == "" {
		c.JSON(400, gin.H{"error": "path parameter is required"})
		return
	}

	// Преобразуем параметры в числа
	var width, height, quality int
	var err error

	if widthStr != "" {
		width, err = strconv.Atoi(widthStr)
		if err != nil || width <= 0 {
			c.JSON(400, gin.H{"error": "invalid width parameter"})
			return
		}
	}

	if heightStr != "" {
		height, err = strconv.Atoi(heightStr)
		if err != nil || height <= 0 {
			c.JSON(400, gin.H{"error": "invalid height parameter"})
			return
		}
	}

	if qualityStr != "" {
		quality, err = strconv.Atoi(qualityStr)
		if err != nil || quality <= 0 || quality > 100 {
			c.JSON(400, gin.H{"error": "invalid quality parameter (1-100)"})
			return
		}
	} else {
		quality = 80 // значение по умолчанию
	}

	// Проверяем, что путь начинается с /static для безопасности
	if !strings.HasPrefix(path, "/static/") {
		c.JSON(400, gin.H{"error": "path must start with /static/"})
		return
	}

	// Убираем начальный слэш для формирования пути к файлу
	filePath := "." + path

	// Получаем оптимизированное изображение
	imgData, err := ih.processor.ProcessImage(c.Request.Context(), filePath, width, height, quality)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to process image: %v", err)})
		return
	}

	// Устанавливаем заголовки
	c.Header("Content-Type", "image/webp")
	c.Header("Cache-Control", "public, max-age=3600") // кэшируем на 1 час

	// Отправляем изображение
	c.Data(200, "image/webp", imgData)
}
