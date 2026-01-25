package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ProcessImage обрабатывает изображение: изменяет размер и конвертирует в WebP
func ProcessImage(filePath string, width, height, quality int) ([]byte, error) {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Проверяем кэш первым
	if cachedData, found := GetCachedImage(filePath, width, height, quality); found {
		return cachedData, nil
	}

	// Открываем исходное изображение
	src, err := openImage(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}

	// Изменяем размер изображения, если указаны параметры
	var dst image.Image
	if width > 0 && height > 0 {
		dst = imaging.Resize(src, width, height, imaging.Lanczos)
	} else if width > 0 {
		dst = imaging.Resize(src, width, 0, imaging.Lanczos)
	} else if height > 0 {
		dst = imaging.Resize(src, 0, height, imaging.Lanczos)
	} else {
		// Если размер не указан, используем оригинальное изображение
		dst = src
	}

	// Конвертируем в оптимизированный формат и возвращаем байты
	buf := new(bytes.Buffer)
	err = encodeOptimizedImage(dst, buf, quality)
	if err != nil {
		return nil, fmt.Errorf("failed to encode optimized image: %v", err)
	}

	result := buf.Bytes()

	// Сохраняем результат в кэш
	SetCachedImage(filePath, width, height, quality, result)

	return result, nil
}

// openImage открывает изображение независимо от его формата
func openImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close() // Игнорируем ошибку закрытия файла
	}()

	// Определяем тип изображения по содержимому
	img, _, err := image.Decode(file)
	return img, err
}

// encodeOptimizedImage кодирует изображение в оптимизированный формат
func encodeOptimizedImage(img image.Image, buf *bytes.Buffer, quality int) error {
	// Используем JPEG с высоким качеством как временный формат
	// В будущем можно добавить поддержку WebP, AVIF и других форматов
	options := &jpeg.Options{Quality: quality}
	return jpeg.Encode(buf, img, options)
}

// GetImageInfo возвращает информацию об изображении
func GetImageInfo(filePath string) (*image.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %v, format: %s", err, format)
	}

	return &config, nil
}

// ValidateImagePath проверяет, что путь к изображению безопасен
func ValidateImagePath(path string) bool {
	// Проверяем, что путь не содержит опасных последовательностей
	normalized := filepath.Clean(path)
	return normalized == path && !filepath.IsAbs(path) && !strings.Contains(normalized, "..")
}
