package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-starter/internal/config"
	"gin-starter/internal/database"
	"gin-starter/internal/handlers"
	"gin-starter/internal/middleware"
	"gin-starter/internal/routes"
	"gin-starter/internal/service/image"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Конфиг
	cfg := config.LoadConfig()

	// 2. Инициализация зависимостей
	image.InitializeCache()
	dbStore, cleanupFunc := database.InitDatabase(cfg)
	// Этот defer сработает только при штатном выходе из main, но для Graceful Shutdown нужно больше
	defer cleanupFunc()

	if dbStore != nil {
		log.Println("✅ Database connection initialized successfully")
	} else {
		log.Println("⚠️ Warning: No database connection established")
	}

	// 3. Роутер
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Статика
	r.StaticFile("/robots.txt", "./static/robots.txt")
	r.StaticFile("/sitemap.xml", "./static/sitemap.xml")
	r.Static("/static", "./static")

	// 4. Сервисы и Хендлеры (DI)
	imageProcessor := image.NewProcessorService()

	// Внедряем dbStore в контекст для доступа в хендлерах
	if dbStore != nil {
		r.Use(func(c *gin.Context) {
			c.Set("dbStore", dbStore)
			c.Next()
		})
	}

	// Создаем обработчики
	pageHandler := handlers.NewPageHandler()
	userHandler := handlers.NewUserHandler()
	imageHandler := handlers.NewImageHandler(imageProcessor)

	// 5. Маршруты
	routes.SetupRoutes(r, pageHandler, userHandler, imageHandler)

	// 6. Запуск сервера с Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Запускаем сервер в горутине, чтобы он не блокировал main
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Ждем сигнала прерывания (Ctrl+C, Docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Даем серверу 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Здесь сработает defer cleanupFunc() перед полным выходом
	log.Println("Server exiting")
}
