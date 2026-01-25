package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	ServerPort string
	DBType     string // "postgres" или "sqlite"
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string // Путь к файлу SQLite
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	// Загружаем .env файл, если он существует
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		DBType:     getEnvOrDefault("DB_TYPE", "sqlite"), // По умолчанию используем SQLite
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:     getEnvOrDefault("DB_NAME", "gin_starter"),
		DBPath:     getEnvOrDefault("DB_PATH", "./data.db"), // Путь к файлу SQLite
	}

	return config
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
