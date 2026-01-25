package image

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheItem представляет элемент кэша
type CacheItem struct {
	Data        []byte
	CreatedAt   time.Time
	AccessCount int
}

// ImageCache простая in-memory кэш-система для изображений
type ImageCache struct {
	items   map[string]*CacheItem
	mutex   sync.RWMutex
	ttl     time.Duration // время жизни кэша
	maxSize int           // максимальный размер кэша в элементах
}

// GlobalCache глобальный экземпляр кэша
var GlobalCache *ImageCache

// InitCache инициализирует глобальный кэш
func InitCache(ttl time.Duration, maxSize int) {
	GlobalCache = &ImageCache{
		items:   make(map[string]*CacheItem),
		ttl:     ttl,
		maxSize: maxSize,
	}

	// Запускаем очистку устаревших элементов
	go GlobalCache.cleanup()
}

// Get получает элемент из кэша
func (c *ImageCache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Проверяем, не истекло ли время жизни
	if time.Since(item.CreatedAt) > c.ttl {
		// Удаляем устаревший элемент
		delete(c.items, key)
		return nil, false
	}

	// Увеличиваем счетчик обращений
	item.AccessCount++
	return item.Data, true
}

// Set сохраняет элемент в кэше
func (c *ImageCache) Set(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если кэш переполнен, удаляем самый старый элемент
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	c.items[key] = &CacheItem{
		Data:      data,
		CreatedAt: time.Now(),
	}
}

// evictOldest удаляет самый старый элемент из кэша
func (c *ImageCache) evictOldest() {
	// В простой реализации удалим первый попавшийся элемент
	// В реальной системе можно использовать LRU или другую стратегию
	for key := range c.items {
		delete(c.items, key)
		break
	}
}

// cleanup периодически удаляет устаревшие элементы
func (c *ImageCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute) // очистка каждые 5 минут
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired удаляет устаревшие элементы
func (c *ImageCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.Sub(item.CreatedAt) > c.ttl {
			delete(c.items, key)
		}
	}
}

// GenerateCacheKey генерирует уникальный ключ для кэширования изображения
func GenerateCacheKey(filePath string, width, height, quality int) string {
	data := fmt.Sprintf("%s_%d_%d_%d", filePath, width, height, quality)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetCachedImage пытается получить изображение из кэша
func GetCachedImage(filePath string, width, height, quality int) ([]byte, bool) {
	if GlobalCache == nil {
		return nil, false
	}

	key := GenerateCacheKey(filePath, width, height, quality)
	return GlobalCache.Get(key)
}

// SetCachedImage сохраняет изображение в кэше
func SetCachedImage(filePath string, width, height, quality int, imageData []byte) {
	if GlobalCache == nil {
		return
	}

	key := GenerateCacheKey(filePath, width, height, quality)
	GlobalCache.Set(key, imageData)
}

// InitializeCache инициализирует кэш при запуске приложения
func InitializeCache() {
	InitCache(1*time.Hour, 100) // TTL 1 час, максимум 100 изображений
}
