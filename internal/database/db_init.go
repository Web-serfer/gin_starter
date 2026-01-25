package database

import (
	"log"

	"gin-starter/internal/config"
	"gin-starter/internal/store"
)

// InitDatabase инициализирует подключение к базе данных в зависимости от типа
func InitDatabase(cfg *config.Config) (store.Store, func()) {
	var dbStore store.Store
	var cleanupFunc func() = func() {} // Пустая функция очистки по умолчанию

	// Создаем подключение к базе данных в зависимости от типа
	switch cfg.DBType {
	case "sqlite":
		sqliteStore, err := store.NewSQLiteStore(cfg.DBPath)
		if err != nil {
			log.Printf("Warning: could not connect to SQLite database: %v", err)
			// Продолжаем работу без базы данных
		} else {
			dbStore = sqliteStore
			// Устанавливаем функцию очистки для закрытия SQLite соединения
			cleanupFunc = func() {
				if err := dbStore.Close(); err != nil {
					log.Printf("Error closing SQLite database: %v", err)
				}
			}
		}
	case "postgres":
		pgStore, err := store.NewPostgreSQLStore(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		if err != nil {
			log.Printf("Warning: could not connect to PostgreSQL database: %v", err)
			// Продолжаем работу без базы данных
		} else {
			dbStore = pgStore
			// Устанавливаем функцию очистки для закрытия PostgreSQL соединения
			cleanupFunc = func() {
				if err := dbStore.Close(); err != nil {
					log.Printf("Error closing PostgreSQL database: %v", err)
				}
			}
		}
	default:
		log.Printf("Warning: unsupported database type: %s", cfg.DBType)
		// Продолжаем работу без базы данных
	}

	// Выполняем миграции
	if dbStore != nil {
		if err := dbStore.Migrate(); err != nil {
			log.Printf("Warning: failed to run migrations: %v", err)
		}
	}

	return dbStore, cleanupFunc
}
