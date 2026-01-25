package routes

import (
	"gin-starter/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// Обратите внимание: я разделил handlers на pageHandler и userApiHandler
func SetupRoutes(r *gin.Engine, pageHandler *handlers.PageHandler, userApiHandler *handlers.UserHandler, imageHandler *handlers.ImageHandler) {

	// 1. Безопасность (через библиотеку надежнее)
	r.Use(secure.New(secure.Config{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	}))

	// 2. CORS (если нужно взаимодействие с внешним фронтендом)
	r.Use(cors.Default())

	// 3. Web-страницы (HTML)
	web := r.Group("/")
	{
		web.GET("/", pageHandler.Home)
		web.GET("/about", pageHandler.About)
		web.GET("/contact", pageHandler.Contact)
		web.GET("/users", pageHandler.Users)
	}

	// 4. Отдельный роут для картинок
	r.GET("/optimized-image", imageHandler.OptimizedImage)

	// 5. API (JSON) с версионированием
	api := r.Group("/api/v1")
	{
		// Тут используем специализированный userApiHandler
		api.GET("/users", userApiHandler.GetUsers)
		api.POST("/users", userApiHandler.CreateUser)
		api.DELETE("/users/:id", userApiHandler.DeleteUser)
	}

	// 6. Обработчик 404 для всех остальных маршрутов
	r.NoRoute(handlers.NotFoundHandler)
}
