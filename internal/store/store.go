package store

import (
	"database/sql"
	"fmt"
	"gin-starter/internal/repository"
)

// Store интерфейс для работы с базой данных
type Store interface {
	Close() error
	Ping() error
	Migrate() error
	// Методы для работы с пользователями
	GetUserRepo() repository.UserRepository
}

// PostgreSQLStore структура для хранения подключений к PostgreSQL
type PostgreSQLStore struct {
	DB       *sql.DB
	UserRepo repository.UserRepository
}

// NewPostgreSQLStore создает новый экземпляр PostgreSQLStore
func NewPostgreSQLStore(dbHost, dbPort, dbUser, dbPassword, dbName string) (*PostgreSQLStore, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &PostgreSQLStore{
		DB: db,
		// Заглушка для репозитория пользователей - будет реализована при необходимости
		UserRepo: nil,
	}

	return store, nil
}

// Close закрывает соединение с базой данных
func (s *PostgreSQLStore) Close() error {
	return s.DB.Close()
}

// Ping проверяет соединение с базой данных
func (s *PostgreSQLStore) Ping() error {
	return s.DB.Ping()
}

// Migrate выполняет миграции (заглушка для PostgreSQL)
func (s *PostgreSQLStore) Migrate() error {
	return nil
}

// GetUserRepo возвращает репозиторий пользователей
func (s *PostgreSQLStore) GetUserRepo() repository.UserRepository {
	// TODO: Реализовать PostgreSQL репозиторий пользователей
	// Пока возвращаем nil, в продакшене нужно будет реализовать
	return s.UserRepo
}
