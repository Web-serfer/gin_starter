package image

import (
	"context"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ProcessorService сервис для обработки изображений
type ProcessorService struct{}

// NewProcessorService создает новый экземпляр сервиса
func NewProcessorService() *ProcessorService {
	return &ProcessorService{}
}

// ProcessImage обрабатывает изображение: изменяет размер и конвертирует в оптимизированный формат
func (ps *ProcessorService) ProcessImage(ctx context.Context, filePath string, width, height, quality int) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	// Открываем исходное изображение
	src, err := ps.openImage(filePath)
	if err != nil {
		return nil, err
	}

	// Изменяем размер изображения, если указаны параметры
	dst := ps.resizeImage(src, width, height)

	// Конвертируем в оптимизированный формат и возвращаем байты
	result, err := ps.encodeOptimizedImage(dst, quality)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// openImage открывает изображение независимо от его формата
func (ps *ProcessorService) openImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	// Определяем тип изображения по содержимому
	img, _, err := image.Decode(file)
	return img, err
}

// resizeImage изменяет размер изображения
func (ps *ProcessorService) resizeImage(src image.Image, width, height int) image.Image {
	if width > 0 && height > 0 {
		return imaging.Resize(src, width, height, imaging.Lanczos)
	} else if width > 0 {
		return imaging.Resize(src, width, 0, imaging.Lanczos)
	} else if height > 0 {
		return imaging.Resize(src, 0, height, imaging.Lanczos)
	} else {
		// Если размер не указан, используем оригинальное изображение
		return src
	}
}

// encodeOptimizedImage кодирует изображение в оптимизированный формат
func (ps *ProcessorService) encodeOptimizedImage(img image.Image, quality int) ([]byte, error) {
	var buf []byte
	writer := &sliceWriter{buf: &buf}

	// Используем JPEG с высоким качеством как временный формат
	// В будущем можно добавить поддержку WebP, AVIF и других форматов
	err := imaging.Encode(writer, img, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// sliceWriter вспомогательная структура для записи в слайс
type sliceWriter struct {
	buf *[]byte
}

func (sw *sliceWriter) Write(p []byte) (n int, err error) {
	*sw.buf = append(*sw.buf, p...)
	return len(p), nil
}

// ValidateImagePath проверяет, что путь к изображению безопасен
func (ps *ProcessorService) ValidateImagePath(path string) bool {
	// Проверяем, что путь не содержит опасных последовательностей
	normalized := filepath.Clean(path)
	return normalized == path && !filepath.IsAbs(path) && !strings.Contains(normalized, "..")
}
