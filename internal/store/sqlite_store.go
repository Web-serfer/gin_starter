package store

import (
	"database/sql"
	"fmt"
	"log"

	"gin-starter/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore структура для хранения подключений к SQLite
type SQLiteStore struct {
	DB       *sql.DB
	UserRepo repository.UserRepository
}

// NewSQLiteStore создает новый экземпляр SQLiteStore
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{
		DB: db,
	}

	// Инициализируем репозитории
	store.UserRepo = repository.NewSQLiteUserRepository(db)

	// Создаем таблицы
	if err := store.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return store, nil
}

// createTables создает необходимые таблицы
func (s *SQLiteStore) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
	}

	for _, query := range queries {
		_, err := s.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func (s *SQLiteStore) Close() error {
	return s.DB.Close()
}

// Ping проверяет соединение с базой данных
func (s *SQLiteStore) Ping() error {
	_, err := s.DB.Exec("SELECT 1")
	return err
}

// Migrate выполняет миграции (в данном случае просто создает таблицы)
func (s *SQLiteStore) Migrate() error {
	log.Println("Running migrations...")

	if err := s.createTables(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// GetUserRepo возвращает репозиторий пользователей
func (s *SQLiteStore) GetUserRepo() repository.UserRepository {
	return s.UserRepo
}
