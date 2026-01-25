package routes

import (
	"gin-starter/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, pageHandler *handlers.PageHandler, imageHandler *handlers.ImageHandler) {

	// --- 1. Middleware Безопасности ---
	// Добавляем этот блок в начало функции, чтобы он применялся ко всем запросам
	r.Use(func(c *gin.Context) {
		// Защита от встраивания в iframe (убирает ошибку X-Frame-Options)
		c.Header("X-Frame-Options", "DENY")

		// Дополнительная защита (рекомендуется)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	})

	// --- 2. Твои маршруты ---
	// (Пример с использованием методов структур, как мы обсуждали)
	web := r.Group("/")
	{
		web.GET("/", pageHandler.Home)
		web.GET("/about", pageHandler.About)
		web.GET("/contact", pageHandler.Contact)
		web.GET("/users", pageHandler.Users)
	}

	r.GET("/optimized-image", imageHandler.OptimizedImage)

	api := r.Group("/api")
	{
		api.GET("/users", pageHandler.GetUsers)
		api.POST("/users", pageHandler.CreateUser)
		api.DELETE("/users/:id", pageHandler.DeleteUser)
		api.POST("/create-test-users", pageHandler.CreateTestUsers)
	}
}
